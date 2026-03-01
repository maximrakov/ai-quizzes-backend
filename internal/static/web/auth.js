const API = '';

// ── Tab switching ──────────────────────────────────────────
const tabs   = document.querySelectorAll('.tab');
const panels = document.querySelectorAll('.panel');

tabs.forEach(tab => {
  tab.addEventListener('click', () => {
    const target = tab.getAttribute('aria-controls');

    tabs.forEach(t => {
      t.classList.toggle('active', t === tab);
      t.setAttribute('aria-selected', String(t === tab));
    });

    panels.forEach(p => {
      p.classList.toggle('hidden', p.id !== target);
    });

    clearMessages();
  });
});

function clearMessages() {
  document.querySelectorAll('.error-msg, .success-msg').forEach(el => {
    el.hidden = true;
    el.textContent = '';
  });
}

// ── Show/hide password ─────────────────────────────────────
document.querySelectorAll('.toggle-pw').forEach(btn => {
  btn.addEventListener('click', () => {
    const input = document.getElementById(btn.dataset.target);
    const isHidden = input.type === 'password';
    input.type = isHidden ? 'text' : 'password';
    btn.setAttribute('aria-label', isHidden ? 'Скрыть пароль' : 'Показать пароль');
  });
});

// ── Helpers ────────────────────────────────────────────────
function setLoading(btn, loading, label) {
  btn.disabled = loading;
  btn.classList.toggle('loading', loading);
  if (!loading) btn.textContent = label;
}

function showError(el, msg) {
  el.textContent = msg;
  el.hidden = false;
}

function showSuccess(el, msg) {
  el.textContent = msg;
  el.hidden = false;
}

function httpErrorMessage(status) {
  switch (status) {
    case 400: return 'Некорректные данные запроса.';
    case 401: return 'Неверное имя пользователя или пароль.';
    case 409: return 'Пользователь с таким именем уже существует.';
    case 500: return 'Внутренняя ошибка сервера. Попробуйте позже.';
    default:  return `Ошибка ${status}. Попробуйте позже.`;
  }
}

// ── Login ──────────────────────────────────────────────────
const formLogin  = document.getElementById('form-login');
const loginError = document.getElementById('login-error');
const loginBtn   = document.getElementById('login-submit');

formLogin.addEventListener('submit', async (e) => {
  e.preventDefault();
  loginError.hidden = true;

  const username = formLogin.username.value.trim();
  const password = formLogin.password.value;

  if (!username || !password) {
    showError(loginError, 'Заполните все поля.');
    return;
  }

  setLoading(loginBtn, true, 'Войти');

  try {
    const res = await fetch(`${API}/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username, password }),
    });

    if (!res.ok) {
      showError(loginError, httpErrorMessage(res.status));
      return;
    }

    const data = await res.json();
    localStorage.setItem('token', data.token);

    // TODO: redirect to dashboard
    window.location.href = 'index.html';
  } catch {
    showError(loginError, 'Не удалось подключиться к серверу.');
  } finally {
    setLoading(loginBtn, false, 'Войти');
  }
});

// ── Register ───────────────────────────────────────────────
const formRegister = document.getElementById('form-register');
const regError     = document.getElementById('reg-error');
const regSuccess   = document.getElementById('reg-success');
const regBtn       = document.getElementById('reg-submit');

formRegister.addEventListener('submit', async (e) => {
  e.preventDefault();
  regError.hidden   = true;
  regSuccess.hidden = true;

  const username = formRegister.username.value.trim();
  const password = formRegister.password.value;
  const role     = formRegister.role.value;

  if (!username || !password) {
    showError(regError, 'Заполните все поля.');
    return;
  }

  if (password.length < 6) {
    showError(regError, 'Пароль должен содержать не менее 6 символов.');
    return;
  }

  setLoading(regBtn, true, 'Создать аккаунт');

  try {
    const res = await fetch(`${API}/user`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username, password, role }),
    });

    if (!res.ok) {
      showError(regError, httpErrorMessage(res.status));
      return;
    }

    formRegister.reset();
    showSuccess(regSuccess, `Аккаунт создан! Теперь вы можете войти.`);

    // auto-switch to login after short delay
    setTimeout(() => {
      document.getElementById('tab-login').click();
      document.getElementById('login-username').value = username;
      document.getElementById('login-username').focus();
    }, 1200);
  } catch {
    showError(regError, 'Не удалось подключиться к серверу.');
  } finally {
    setLoading(regBtn, false, 'Создать аккаунт');
  }
});
