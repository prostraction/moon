import { initializeApp } from './ui.js';
import { createStars, createClouds } from './render.js';

document.addEventListener("DOMContentLoaded", () => {
  createStars();
  createClouds();
  initializeApp();
});
