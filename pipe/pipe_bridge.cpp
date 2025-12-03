#include <cstring>
#include <fcntl.h>
#include <iostream>
#include <poll.h>
#include <signal.h>
#include <string>
#include <sys/stat.h>
#include <unistd.h>

volatile bool running = true;

void signal_handler(int signal) {
  if (signal == SIGINT || signal == SIGTERM) {
    running = false;
  }
}

int main(int argc, char *argv[]) {
  if (argc != 3) {
    std::cerr << "Usage: " << argv[0] << " <pipe1> <pipe2>" << std::endl;
    std::cerr << "  pipe1: stdin will be written to this pipe" << std::endl;
    std::cerr << "  pipe2: data from this pipe will be written to stdout"
              << std::endl;
    return 1;
  }

  std::string pipe1_name = argv[1];
  std::string pipe2_name = argv[2];

  // Create pipes if they don't exist
  if (mkfifo(pipe1_name.c_str(), 0666) == -1 && errno != EEXIST) {
    std::cerr << "Error: Cannot create pipe1 '" << pipe1_name
              << "': " << strerror(errno) << std::endl;
    return 1;
  }

  if (mkfifo(pipe2_name.c_str(), 0666) == -1 && errno != EEXIST) {
    std::cerr << "Error: Cannot create pipe2 '" << pipe2_name
              << "': " << strerror(errno) << std::endl;
    return 1;
  }

  // Set up signal handlers
  signal(SIGINT, signal_handler);
  signal(SIGTERM, signal_handler);

  // Open pipe1 for writing (stdin -> pipe1)
  // First try non-blocking, then fall back to blocking if needed
  int pipe1_fd = open(pipe1_name.c_str(), O_WRONLY | O_NONBLOCK);
  if (pipe1_fd == -1) {
    if (errno == ENXIO) {
      // No readers yet, open in blocking mode
      std::cerr << "Waiting for reader on " << pipe1_name << "..." << std::endl;
      pipe1_fd = open(pipe1_name.c_str(), O_WRONLY);
    } else {
      std::cerr << "Error: Cannot open pipe1 '" << pipe1_name
                << "' for writing: " << strerror(errno) << std::endl;
      return 1;
    }
  }
  if (pipe1_fd == -1) {
    std::cerr << "Error: Cannot open pipe1 '" << pipe1_name
              << "' for writing: " << strerror(errno) << std::endl;
    return 1;
  }

  // Open pipe2 for reading (pipe2 -> stdout)
  int pipe2_fd = open(pipe2_name.c_str(), O_RDONLY | O_NONBLOCK);
  if (pipe2_fd == -1) {
    std::cerr << "Error: Cannot open pipe2 '" << pipe2_name
              << "' for reading: " << strerror(errno) << std::endl;
    close(pipe1_fd);
    return 1;
  }

  // Set stdin to non-blocking mode
  int flags = fcntl(STDIN_FILENO, F_GETFL, 0);
  fcntl(STDIN_FILENO, F_SETFL, flags | O_NONBLOCK);

  // Set stdout to non-blocking mode
  flags = fcntl(STDOUT_FILENO, F_GETFL, 0);
  fcntl(STDOUT_FILENO, F_SETFL, flags | O_NONBLOCK);

  const int BUFFER_SIZE = 4096;
  char buffer[BUFFER_SIZE];

  std::cerr << "Pipe bridge started. stdin -> " << pipe1_name << ", "
            << pipe2_name << " -> stdout" << std::endl;
  std::cerr << "Press Ctrl+C to stop." << std::endl;

  while (running) {
    // Set up poll structure
    struct pollfd fds[3];
    fds[0].fd = STDIN_FILENO;
    fds[0].events = POLLIN;
    fds[0].revents = 0;

    fds[1].fd = pipe2_fd;
    fds[1].events = POLLIN;
    fds[1].revents = 0;

    fds[2].fd = STDOUT_FILENO;
    fds[2].events = POLLOUT;
    fds[2].revents = 0;

    // Wait for activity
    int result = poll(fds, 3, 100); // 100ms timeout
    if (result == -1) {
      if (errno == EINTR) {
        continue; // Interrupted by signal
      }
      std::cerr << "Error in poll: " << strerror(errno) << std::endl;
      break;
    }

    // Handle stdin -> pipe1
    if (fds[0].revents & POLLIN) {
      ssize_t bytes_read = read(STDIN_FILENO, buffer, BUFFER_SIZE);
      if (bytes_read > 0) {
        ssize_t bytes_written = write(pipe1_fd, buffer, bytes_read);
        if (bytes_written == -1) {
          if (errno != EAGAIN && errno != EWOULDBLOCK) {
            std::cerr << "Error writing to pipe1: " << strerror(errno)
                      << std::endl;
            break;
          }
        }
      } else if (bytes_read == 0) {
        // EOF on stdin
        running = false;
      }
    }

    // Handle pipe2 -> stdout
    if (fds[1].revents & POLLIN) {
      ssize_t bytes_read = read(pipe2_fd, buffer, BUFFER_SIZE);
      if (bytes_read > 0) {
        ssize_t bytes_written = write(STDOUT_FILENO, buffer, bytes_read);
        if (bytes_written == -1) {
          if (errno != EAGAIN && errno != EWOULDBLOCK) {
            std::cerr << "Error writing to stdout: " << strerror(errno)
                      << std::endl;
            break;
          }
        }
      } else if (bytes_read == 0) {
        // Pipe closed
        std::cerr << "Warning: pipe2 closed" << std::endl;
        running = false;
      }
    }

    // Check for errors
    if (fds[0].revents & (POLLERR | POLLHUP)) {
      std::cerr << "Error on stdin" << std::endl;
      break;
    }
    if (fds[1].revents & (POLLERR | POLLHUP)) {
      std::cerr << "Error on pipe2" << std::endl;
      break;
    }
  }

  // Clean up
  close(pipe1_fd);
  close(pipe2_fd);

  std::cerr << "\nPipe bridge stopped." << std::endl;
  return 0;
}
