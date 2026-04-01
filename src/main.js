const API_BASE = 'http://localhost:8080'

document.addEventListener('DOMContentLoaded', () => {
  loadRecipes();
  setupEventListeners();
});

function setupEventListeners() {
  document.getElementById('add-recipe-btn').addEventListener('click', () => openModal());
  document.getElementById('search-btn').addEventListener('click', searchRecipes);
  document.getElementById('search-input').addEventListener('keypress', (e) => {
    if (e.key === 'Enter') searchRecipes();
  });
  document.getElementById('recipe-form').addEventListener('submit', handleFormSubmit);
  document.querySelector('.close').addEventListener('click', closeModal);
}

async function loadRecipes() {
  try {
    const response = await fetch(`${API_BASE}/recipes`);
    if (!response.ok) {
      console.error('Failed to load recipes:', response.status);
      return;
    }
    const recipes = await response.json();
    displayRecipes(recipes);
  } catch (error) {
    console.error('Failed to load recipes:', error);
  }
}

async function searchRecipes() {
  const query = document.getElementById('search-input').value;
  const url = query ? `${API_BASE}/recipes/search?q=${encodeURIComponent(query)}` : `${API_BASE}/recipes`;

  try {
    const response = await fetch(url);
    const recipes = await response.json();
    displayRecipes(recipes);
  } catch (error) {
    console.error('Search failed:', error);
    displayRecipes([]);
  }
}

function displayRecipes(recipes) {
  const grid = document.getElementById('recipe-grid');
  if (!grid) {
    console.error('recipe-grid not found');
    return;
  }

  if (!recipes || recipes.length === 0) {
    grid.innerHTML = '<p>No recipes found.</p>';
    return;
  }

  grid.innerHTML = recipes.map(recipe => `
    <div class="recipe-card">
      <h3>${escapeHtml(recipe.title)}</h3>
      <div class="cuisine">${escapeHtml(recipe.cuisine || 'No cuisine')}</div>
      <div class="meta">
        ${recipe.cooking_time || 0} min | ${recipe.servings || 0} servings      
      </div>
      <div class="ingredients">
        ${Array.isArray(recipe.ingredients) ? recipe.ingredients.join(', ') : ''}
      </div>
      <div class="actions">
        <button class="btn-edit" onclick="editRecipe(${recipe.id})">Edit</button>
        <button class="btn-delete" onclick="deleteRecipe(${recipe.id})">Delete</button>
      </div>
    </div>
  `).join('');
}

function openModal(recipe = null) {
  const modal = document.getElementById('recipe-modal');
  const form = document.getElementById('recipe-form');
  const modalTitle = document.getElementById('modal-title');

  form.reset();
  if (recipe) {
    modalTitle.textContent = 'Edit Recipe';
    document.getElementById('recipe-id').value = recipe.id;
    document.getElementById('title').value = recipe.title || '';
    document.getElementById('url').value = recipe.url || '';
    document.getElementById('cuisine').value = recipe.cuisine || '';
    document.getElementById('cooking_time').value = recipe.cooking_time || '';
    document.getElementById('servings').value = recipe.servings || '';
    document.getElementById('ingredients').value = Array.isArray(recipe.ingredients) ? recipe.ingredients.join(', ') : '';
    document.getElementById('instructions').value = recipe.instructions || '';
  } else {
    modalTitle.textContent = 'Add Recipe';
    document.getElementById('recipe-id').value = '';
  }
  modal.style.display = 'block';
}

function closeModal() {
  document.getElementById('recipe-modal').style.display = 'none';
}

async function handleFormSubmit(e) {
  e.preventDefault();

  const id = document.getElementById('recipe-id').value;
  const ingredientsText = document.getElementById('ingredients').value;
  const ingredients = ingredientsText ? ingredientsText.split(',').map(s => s.trim()) : [];

  const recipe = {
    title: document.getElementById('title').value,
    url: document.getElementById('url').value,
    cuisine: document.getElementById('cuisine').value,
    cooking_time: parseInt(document.getElementById('cooking_time').value) || 0,
    servings: parseInt(document.getElementById('servings').value) || 0,
    ingredients: ingredients,
    instructions: document.getElementById('instructions').value
  };

  try {
    const url = id ? `${API_BASE}/recipes/${id}` : `${API_BASE}/recipes`;
    const method = id ? 'PUT' : 'POST';

    const response = await fetch(url, {
      method: method,
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify(recipe)
    });

    if (response.ok) {
      closeModal();
      loadRecipes();
    }
  } catch (error) {
    console.error('Failed to save recipe:', error);
  }
}

async function editRecipe(id) {
  try {
    const response = await fetch(`${API_BASE}/recipes/${id}`);
    const recipe = await response.json();
    openModal(recipe);
  } catch (error) {
    console.error('Failed to load recipe:', error);
  }
}

async function deleteRecipe(id) {
  if (!confirm('Are you sure you want to delete this recipe?')) return;

  try {
    await fetch(`${API_BASE}/recipes/${id}`, {method: 'DELETE'});
    loadRecipes();
  } catch (error) {
    console.error('Failed to delete recipe:', error);
  }
}

function escapeHtml(text) {
  if (!text) return '';
  const div = document.createElement('div');
  div.textContent = text;
  return div.innerHTML;
}
