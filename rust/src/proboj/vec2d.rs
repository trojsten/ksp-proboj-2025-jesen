use derive_more::{Add, Sub};
use serde::{Deserialize, Serialize};
use std::ops::Mul;

#[derive(Clone, Copy, Debug, Serialize, Deserialize, Add, Sub)]
pub struct Vec2D {
    pub x: f64,
    pub y: f64,
}

impl Vec2D {
    pub fn distance(&self, other: &Vec2D) -> f64 {
        ((self.x - other.x).powi(2) + (self.y - other.y).powi(2)).sqrt()
    }

    pub fn magnitude(&self) -> f64 {
        (self.x.powi(2) + self.y.powi(2)).sqrt()
    }

    pub fn normalize(&self) -> Vec2D {
        let mag = self.magnitude();
        Vec2D {
            x: self.x / mag,
            y: self.y / mag,
        }
    }

    pub fn rotate(&self, angle_rad: f64) -> Vec2D {
        let cos_theta = angle_rad.cos();
        let sin_theta = angle_rad.sin();
        Vec2D {
            x: self.x * cos_theta - self.y * sin_theta,
            y: self.x * sin_theta + self.y * cos_theta,
        }
    }

    pub fn clamp_magnitude(&self, max_magnitude: f64) -> Vec2D {
        let mag = self.magnitude();
        if mag > max_magnitude {
            self.normalize() * max_magnitude
        } else {
            *self
        }
    }

    pub fn zero() -> Vec2D {
        Vec2D { x: 0.0, y: 0.0 }
    }
}

impl Mul<f64> for Vec2D {
    type Output = Vec2D;

    fn mul(self, rhs: f64) -> Self::Output {
        Vec2D {
            x: self.x * rhs,
            y: self.y * rhs,
        }
    }
}

impl Mul<Vec2D> for f64 {
    type Output = Vec2D;

    fn mul(self, rhs: Vec2D) -> Self::Output {
        Vec2D {
            x: rhs.x * self,
            y: rhs.y * self,
        }
    }
}
