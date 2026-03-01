// ── State ──────────────────────────────────────────────────
let token, userId, role, username;
let questionCount = 0;

// ── Bootstrap ──────────────────────────────────────────────
function init() {
  token = localStorage.getItem('token');
  if (!token) { redirect('auth.html'); return; }

  let claims;
  try {
    claims = parseJwt(token);
  } catch {
    redirect('auth.html');
    return;
  }

  userId   = claims.user_id;
  role     = claims.role;
  username = claims.sub || `User #${userId}`;

  document.getElementById('topbar-username').textContent = username;
  document.getElementById('role-badge').textContent = role === 'mentor' ? 'Ментор' : 'Студент';
  document.getElementById('logout-btn').addEventListener('click', logout);

  renderNav();
  navigate('my-quizzes');
}

// ── Auth helpers ───────────────────────────────────────────
function parseJwt(t) {
  const b64 = t.split('.')[1].replace(/-/g, '+').replace(/_/g, '/');
  return JSON.parse(atob(b64));
}

function logout() {
  localStorage.clear();
  redirect('auth.html');
}

function redirect(path) {
  window.location.href = path;
}

// ── API helper ─────────────────────────────────────────────
async function api(method, path, body) {
  const opts = {
    method,
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`,
    },
  };
  if (body !== undefined) opts.body = JSON.stringify(body);
  return fetch(path, opts);
}

// ── Navigation ─────────────────────────────────────────────
const NAV = {
  mentor: [
    { id: 'my-quizzes',   label: 'Мои квизы',    icon: navIconQuizzes() },
    { id: 'create-quiz',  label: 'Создать квиз',  icon: navIconCreate()  },
  ],
  student: [
    { id: 'my-quizzes',  label: 'Мои квизы',      icon: navIconQuizzes() },
    { id: 'results',     label: 'Мои результаты',  icon: navIconResults() },
  ],
};

function renderNav() {
  const items = NAV[role] ?? NAV.student;
  const sidebar = document.getElementById('sidebar');
  sidebar.innerHTML = items.map(item => `
    <button class="nav-item" data-view="${item.id}" onclick="navigate('${item.id}')">
      ${item.icon}
      <span>${item.label}</span>
    </button>
  `).join('');
}

function setActiveNav(viewId) {
  document.querySelectorAll('.nav-item').forEach(btn => {
    btn.classList.toggle('active', btn.dataset.view === viewId);
  });
}

async function navigate(viewId) {
  setActiveNav(viewId);
  const main = document.getElementById('main');
  main.innerHTML = '<div class="page-loading"><div class="spinner"></div></div>';

  switch (viewId) {
    case 'my-quizzes':
      if (role === 'mentor') await renderMentorQuizzes();
      else                   await renderStudentQuizzes();
      break;
    case 'create-quiz':
      renderCreateQuiz();
      break;
    case 'results':
      renderResults();
      break;
  }
}

// ── Mentor: my quizzes ─────────────────────────────────────
async function renderMentorQuizzes() {
  const main = document.getElementById('main');
  try {
    const res = await api('GET', `/user/${userId}/quizzes/created`);
    if (!res.ok) throw new Error(res.status);
    const quizzes = await res.json();

    if (!quizzes || quizzes.length === 0) {
      main.innerHTML = pageHeader('Мои квизы') + `
        <div class="empty-state">
          <p>У вас пока нет созданных квизов.</p>
          <button class="btn-primary" onclick="navigate('create-quiz')">Создать первый квиз</button>
        </div>`;
      return;
    }

    main.innerHTML = pageHeader('Мои квизы', quizzes.length) + `
      <div class="quiz-list" id="quiz-list"></div>`;

    document.getElementById('quiz-list').innerHTML = quizzes.map(quiz => `
      <div class="quiz-card">
        <div class="quiz-card-header">
          <div class="quiz-card-info">
            <h3>${esc(quiz.title)}</h3>
            <span class="quiz-meta">${quiz.questions?.length ?? 0} вопросов · ID ${quiz.id}</span>
          </div>
          <button class="btn-outline btn-sm" onclick="openAssignModal(${quiz.id}, '${esc(quiz.title)}')">
            Назначить
          </button>
        </div>
        <details class="quiz-accordion">
          <summary>Показать вопросы</summary>
          <div class="accordion-body">
            ${(quiz.questions ?? []).map((q, i) => `
              <div class="question-view">
                <p class="question-view-text">${i + 1}. ${esc(q.text)}</p>
                <ul class="options-view">
                  ${(q.options ?? []).sort((a, b) => a.number - b.number).map(opt => `
                    <li class="${opt.number === q.correct_answer_number ? 'opt-correct' : 'opt-normal'}">
                      ${opt.number}. ${esc(opt.text)}
                      ${opt.number === q.correct_answer_number ? '<span class="checkmark">✓</span>' : ''}
                    </li>`).join('')}
                </ul>
              </div>`).join('')}
          </div>
        </details>
      </div>`).join('');

  } catch (e) {
    main.innerHTML = errorState(`Не удалось загрузить квизы: ${e.message}`);
  }
}

// ── Student: my quizzes ────────────────────────────────────
async function renderStudentQuizzes() {
  const main = document.getElementById('main');
  try {
    const [assignRes, quizRes] = await Promise.all([
      api('GET', `/user/${userId}/assignments`),
      api('GET', `/user/${userId}/quizzes/assigned`),
    ]);

    if (!assignRes.ok || !quizRes.ok) throw new Error('Ошибка загрузки');

    const assignments = await assignRes.json() ?? [];
    const quizzes     = await quizRes.json()     ?? [];

    const quizById = {};
    quizzes.forEach(q => { quizById[q.id] = q; });

    const done = getCompleted();

    if (assignments.length === 0) {
      main.innerHTML = pageHeader('Мои квизы') + `
        <div class="empty-state"><p>Вам пока не назначено ни одного квиза.</p></div>`;
      return;
    }

    main.innerHTML = pageHeader('Мои квизы', assignments.length) + `
      <div class="quiz-list" id="quiz-list"></div>`;

    document.getElementById('quiz-list').innerHTML = assignments.map(a => {
      const quiz    = quizById[a.quiz_id];
      const isDone  = done.includes(a.id);
      return `
        <div class="quiz-card ${isDone ? 'quiz-card-done' : ''}">
          <div class="quiz-card-header">
            <div class="quiz-card-info">
              <h3>${quiz ? esc(quiz.title) : `Квиз #${a.quiz_id}`}</h3>
              ${quiz ? `<span class="quiz-meta">${quiz.questions?.length ?? 0} вопросов</span>` : ''}
            </div>
            ${isDone
              ? '<span class="status-badge status-done">Пройден</span>'
              : `<button class="btn-primary btn-sm"
                   onclick="startQuiz(${a.id}, ${a.quiz_id})">Пройти</button>`}
          </div>
        </div>`;
    }).join('');

  } catch (e) {
    main.innerHTML = errorState(`Не удалось загрузить квизы: ${e.message}`);
  }
}

// ── Quiz taking ────────────────────────────────────────────
async function startQuiz(assignmentId, quizId) {
  const overlay    = document.getElementById('quiz-overlay');
  const container  = document.getElementById('quiz-container');
  overlay.classList.remove('hidden');
  container.innerHTML = '<div class="page-loading"><div class="spinner"></div></div>';

  try {
    const res = await api('GET', `/quiz/${quizId}`);
    if (!res.ok) throw new Error(`HTTP ${res.status}`);
    renderQuizTaking(assignmentId, await res.json());
  } catch (e) {
    container.innerHTML = `
      <div class="overlay-header">
        <button class="btn-ghost" onclick="closeQuizOverlay()">← Назад</button>
      </div>
      ${errorState(`Не удалось загрузить квиз: ${e.message}`)}`;
  }
}

function renderQuizTaking(assignmentId, quiz) {
  const container = document.getElementById('quiz-container');
  container.innerHTML = `
    <div class="overlay-header">
      <button class="btn-ghost" onclick="closeQuizOverlay()">← Назад</button>
      <h2 class="overlay-title">${esc(quiz.title)}</h2>
      <span class="quiz-meta">${quiz.questions.length} вопросов</span>
    </div>
    <form id="quiz-form" class="quiz-form">
      ${quiz.questions.map((q, i) => `
        <div class="quiz-question">
          <p class="question-view-text">${i + 1}. ${esc(q.text)}</p>
          <div class="options-group">
            ${q.options.sort((a, b) => a.number - b.number).map(opt => `
              <label class="option-label">
                <input type="radio" name="q${q.id}" value="${opt.number}" required />
                <span>${opt.number}. ${esc(opt.text)}</span>
              </label>`).join('')}
          </div>
        </div>`).join('')}
      <div id="quiz-submit-error" class="error-msg" hidden></div>
      <div class="quiz-submit-row">
        <button type="submit" class="btn-primary" id="submit-btn">Завершить</button>
      </div>
    </form>`;

  document.getElementById('quiz-form').addEventListener('submit', e => {
    e.preventDefault();
    submitQuiz(assignmentId, quiz);
  });
}

async function submitQuiz(assignmentId, quiz) {
  const form    = document.getElementById('quiz-form');
  const btn     = document.getElementById('submit-btn');
  const errorEl = document.getElementById('quiz-submit-error');
  errorEl.hidden = true;

  const answers = quiz.questions.map(q => ({
    question_id:   q.id,
    answer_number: parseInt(form.querySelector(`input[name="q${q.id}"]:checked`).value),
  }));

  btn.disabled = true;
  btn.classList.add('loading');
  btn.textContent = 'Отправка';

  try {
    const res = await api('POST', '/attempt', { assignment_id: assignmentId, answers });
    if (!res.ok) throw new Error(`HTTP ${res.status}`);
    const attempt = await res.json();

    saveResult(assignmentId, quiz.title, attempt);
    markCompleted(assignmentId);
    closeQuizOverlay();
    showResultModal(quiz, attempt);
  } catch (e) {
    errorEl.textContent = `Ошибка отправки: ${e.message}`;
    errorEl.hidden = false;
    btn.disabled = false;
    btn.classList.remove('loading');
    btn.textContent = 'Завершить';
  }
}

function closeQuizOverlay() {
  document.getElementById('quiz-overlay').classList.add('hidden');
}

// ── Result modal ───────────────────────────────────────────
function showResultModal(quiz, attempt) {
  const backdrop = document.getElementById('modal-backdrop');
  const box      = document.getElementById('modal-box');

  const correctSet = new Set(attempt.correct_question_ids ?? []);
  const score      = attempt.score ?? 0;
  const cls        = score >= 70 ? 'score-good' : score >= 40 ? 'score-mid' : 'score-bad';

  box.innerHTML = `
    <h2 class="modal-title">Результат квиза</h2>
    <div class="score-circle ${cls}">${score.toFixed(0)}%</div>
    <p class="score-label">${scoreLabel(score)}</p>
    <div class="result-list">
      ${quiz.questions.map(q => {
        const ok = correctSet.has(q.id);
        return `<div class="result-row ${ok ? 'result-ok' : 'result-fail'}">
          <span class="result-icon">${ok ? '✓' : '✗'}</span>
          <span>${esc(q.text)}</span>
        </div>`;
      }).join('')}
    </div>
    <button class="btn-primary" onclick="closeModal()">Закрыть</button>`;

  backdrop.classList.remove('hidden');
}

function scoreLabel(s) {
  if (s >= 90) return 'Отлично!';
  if (s >= 70) return 'Хорошо!';
  if (s >= 50) return 'Неплохо';
  return 'Стоит повторить тему';
}

function closeModal() {
  document.getElementById('modal-backdrop').classList.add('hidden');
  navigate('my-quizzes');
}

// ── Results view ───────────────────────────────────────────
function renderResults() {
  const main    = document.getElementById('main');
  const results = getSavedResults();

  if (results.length === 0) {
    main.innerHTML = pageHeader('Мои результаты') + `
      <div class="empty-state"><p>Вы ещё не проходили ни одного квиза.</p></div>`;
    return;
  }

  main.innerHTML = pageHeader('Мои результаты', results.length) + `
    <div class="quiz-list">
      ${[...results].reverse().map(r => {
        const cls = r.score >= 70 ? 'score-good' : r.score >= 40 ? 'score-mid' : 'score-bad';
        return `
          <div class="quiz-card">
            <div class="quiz-card-header">
              <div class="quiz-card-info">
                <h3>${esc(r.quizTitle)}</h3>
                <span class="quiz-meta">${new Date(r.date).toLocaleDateString('ru-RU')}</span>
              </div>
              <div class="score-pill ${cls}">${r.score.toFixed(0)}%</div>
            </div>
            <div class="result-summary">
              <span class="summary-ok">✓ ${r.correctCount} правильно</span>
              <span class="summary-fail">✗ ${r.wrongCount} неправильно</span>
            </div>
          </div>`;
      }).join('')}
    </div>`;
}

// ── Create quiz ────────────────────────────────────────────
function renderCreateQuiz() {
  const main = document.getElementById('main');
  main.innerHTML = pageHeader('Создать квиз') + `
    <div class="tabs" role="tablist">
      <button class="tab active" id="tab-manual" onclick="switchCreateTab('manual')">Вручную</button>
      <button class="tab"        id="tab-ai"     onclick="switchCreateTab('ai')">Генерация AI</button>
    </div>
    <div id="create-panel"></div>`;
  renderManualForm();
}

function switchCreateTab(tab) {
  document.getElementById('tab-manual').classList.toggle('active', tab === 'manual');
  document.getElementById('tab-ai').classList.toggle('active', tab === 'ai');
  if (tab === 'manual') renderManualForm();
  else                  renderAiForm();
}

// Manual form
function renderManualForm() {
  const panel = document.getElementById('create-panel');
  questionCount = 0;
  panel.innerHTML = `
    <form id="manual-form" class="create-form" onsubmit="submitManualQuiz(event)">
      <div class="field">
        <label for="quiz-title">Название квиза</label>
        <input id="quiz-title" type="text" placeholder="Основы Go" required />
      </div>
      <div id="questions-container"></div>
      <button type="button" class="btn-outline" onclick="addQuestion()">+ Добавить вопрос</button>
      <div id="create-error"   class="error-msg"   hidden></div>
      <div id="create-success" class="success-msg"  hidden></div>
      <button type="submit" class="btn-primary" id="create-btn">Создать квиз</button>
    </form>`;
  addQuestion();
}

function addQuestion() {
  questionCount++;
  const n   = questionCount;
  const div = document.createElement('div');
  div.className = 'question-builder';
  div.id = `qb-${n}`;
  div.innerHTML = `
    <div class="qb-header">
      <span class="qb-label">Вопрос ${n}</span>
      ${n > 1 ? `<button type="button" class="btn-danger-ghost" onclick="removeQuestion(${n})">Удалить</button>` : ''}
    </div>
    <div class="field">
      <label>Текст вопроса</label>
      <input type="text" name="q${n}_text" placeholder="Что такое горутина в Go?" required />
    </div>
    <div class="field"><label>Варианты ответа <span class="hint-inline">(отметьте правильный)</span></label></div>
    <div class="options-builder">
      ${[1, 2, 3, 4].map(i => `
        <div class="option-row">
          <input type="radio" name="q${n}_correct" value="${i}" ${i === 1 ? 'checked' : ''} />
          <input type="text" name="q${n}_opt${i}" placeholder="Вариант ${i}" required />
        </div>`).join('')}
    </div>`;
  document.getElementById('questions-container').appendChild(div);
}

function removeQuestion(n) {
  document.getElementById(`qb-${n}`)?.remove();
}

async function submitManualQuiz(e) {
  e.preventDefault();
  const form      = e.target;
  const errorEl   = document.getElementById('create-error');
  const successEl = document.getElementById('create-success');
  const btn       = document.getElementById('create-btn');
  errorEl.hidden = successEl.hidden = true;

  const title     = document.getElementById('quiz-title').value.trim();
  const questions = [];

  document.querySelectorAll('.question-builder').forEach(qb => {
    const n       = qb.id.replace('qb-', '');
    const text    = form.querySelector(`[name="q${n}_text"]`)?.value.trim();
    const correct = parseInt(form.querySelector(`[name="q${n}_correct"]:checked`)?.value);
    const options = [1, 2, 3, 4].map(i => ({
      number: i,
      text: form.querySelector(`[name="q${n}_opt${i}"]`)?.value.trim() ?? '',
    }));
    if (text) questions.push({ text, correct_answer_number: correct, options });
  });

  if (questions.length === 0) {
    errorEl.textContent = 'Добавьте хотя бы один вопрос.';
    errorEl.hidden = false;
    return;
  }

  btn.disabled = true;
  btn.classList.add('loading');
  btn.textContent = 'Создание';

  try {
    const res = await api('POST', '/quiz', { title, questions });
    if (!res.ok) throw new Error(res.status === 500 ? 'Только менторы могут создавать квизы.' : `Ошибка ${res.status}`);
    const quiz = await res.json();
    successEl.textContent = `Квиз «${quiz.title}» успешно создан!`;
    successEl.hidden = false;
    questionCount = 0;
    document.getElementById('questions-container').innerHTML = '';
    document.getElementById('quiz-title').value = '';
    addQuestion();
  } catch (err) {
    errorEl.textContent = err.message;
    errorEl.hidden = false;
  } finally {
    btn.disabled = false;
    btn.classList.remove('loading');
    btn.textContent = 'Создать квиз';
  }
}

// AI form
function renderAiForm() {
  document.getElementById('create-panel').innerHTML = `
    <form id="ai-form" class="create-form" onsubmit="submitAiQuiz(event)">
      <div class="field">
        <label for="ai-title">Название квиза</label>
        <input id="ai-title" type="text" placeholder="Квиз по основам Go" required />
      </div>
      <div class="field">
        <label for="ai-topic">Тема</label>
        <input id="ai-topic" type="text" placeholder="Основы Go" required />
      </div>
      <div class="field">
        <label for="ai-count">Количество вопросов</label>
        <input id="ai-count" type="number" value="5" min="1" max="20" required />
      </div>
      <div id="ai-error" class="error-msg" hidden></div>
      <button type="submit" class="btn-primary" id="ai-btn">Сгенерировать</button>
    </form>
    <div id="ai-result" hidden></div>`;
}

async function submitAiQuiz(e) {
  e.preventDefault();
  const errorEl  = document.getElementById('ai-error');
  const resultEl = document.getElementById('ai-result');
  const btn      = document.getElementById('ai-btn');
  errorEl.hidden = true;
  resultEl.hidden = true;

  const title = document.getElementById('ai-title').value.trim();
  const topic = document.getElementById('ai-topic').value.trim();
  const count = parseInt(document.getElementById('ai-count').value);

  btn.disabled = true;
  btn.classList.add('loading');
  btn.textContent = 'Генерация';

  try {
    const res = await api('POST', '/quiz/generate', { title, topic, count });
    if (!res.ok) throw new Error(`Ошибка ${res.status}`);
    const quiz = await res.json();
    renderGeneratedQuiz(quiz);
  } catch (err) {
    errorEl.textContent = err.message;
    errorEl.hidden = false;
  } finally {
    btn.disabled = false;
    btn.classList.remove('loading');
    btn.textContent = 'Сгенерировать снова';
  }
}

// ── Assign modal ───────────────────────────────────────────
async function openAssignModal(quizId, quizTitle) {
  const backdrop = document.getElementById('assign-backdrop');
  const box      = document.getElementById('assign-box');

  box.innerHTML = `
    <h2 class="modal-title">Назначить квиз</h2>
    <p class="assign-quiz-name">${esc(quizTitle)}</p>
    <div class="page-loading"><div class="spinner"></div></div>`;
  backdrop.classList.remove('hidden');

  try {
    const [studentsRes, assignmentsRes] = await Promise.all([
      api('GET', '/user?role=student'),
      api('GET', `/quiz/${quizId}/assignments`),
    ]);

    if (!studentsRes.ok || !assignmentsRes.ok) throw new Error('Ошибка загрузки данных');

    const students    = await studentsRes.json()    ?? [];
    const assignments = await assignmentsRes.json() ?? [];

    const assignedIds = new Set(assignments.map(a => a.student_id));
    const available   = students.filter(s => !assignedIds.has(s.id));

    renderAssignForm(box, quizId, quizTitle, available);
  } catch (e) {
    box.innerHTML += `<div class="error-state" style="margin-top:1rem">${esc(e.message)}</div>
      <div class="assign-actions" style="margin-top:1rem">
        <button class="btn-ghost" onclick="closeAssignModal()">Закрыть</button>
      </div>`;
  }
}

function renderAssignForm(box, quizId, quizTitle, students) {
  if (students.length === 0) {
    box.innerHTML = `
      <h2 class="modal-title">Назначить квиз</h2>
      <p class="assign-quiz-name">${esc(quizTitle)}</p>
      <p class="assign-empty">Все студенты уже получили этот квиз.</p>
      <div class="assign-actions">
        <button class="btn-ghost" onclick="closeAssignModal()">Закрыть</button>
      </div>`;
    return;
  }

  box.innerHTML = `
    <h2 class="modal-title">Назначить квиз</h2>
    <p class="assign-quiz-name">${esc(quizTitle)}</p>
    <p class="assign-hint">Выберите одного или нескольких студентов:</p>
    <div class="student-list" id="student-list">
      ${students.map(s => `
        <label class="student-row">
          <input type="checkbox" class="student-cb" value="${s.id}" />
          <span class="student-name">${esc(s.username)}</span>
          <span class="student-id">ID ${s.id}</span>
        </label>`).join('')}
    </div>
    <div id="assign-error"   class="error-msg"  hidden></div>
    <div id="assign-success" class="success-msg" hidden></div>
    <div class="assign-actions">
      <button class="btn-ghost" onclick="closeAssignModal()">Отмена</button>
      <button class="btn-primary" id="assign-btn"
              style="width:auto;margin-top:0"
              onclick="submitAssignment(${quizId})">Назначить</button>
    </div>`;
}

function closeAssignModal() {
  document.getElementById('assign-backdrop').classList.add('hidden');
}

async function submitAssignment(quizId) {
  const errorEl   = document.getElementById('assign-error');
  const successEl = document.getElementById('assign-success');
  const btn       = document.getElementById('assign-btn');
  errorEl.hidden = successEl.hidden = true;

  const selected = [...document.querySelectorAll('.student-cb:checked')].map(cb => parseInt(cb.value));
  if (selected.length === 0) {
    errorEl.textContent = 'Выберите хотя бы одного студента.';
    errorEl.hidden = false;
    return;
  }

  btn.disabled = true;
  btn.classList.add('loading');
  btn.textContent = 'Назначение';

  const failed = [];
  for (const studentId of selected) {
    const res = await api('POST', '/assignment', { quiz_id: quizId, student_id: studentId });
    if (!res.ok) failed.push(studentId);
  }

  btn.disabled = false;
  btn.classList.remove('loading');
  btn.textContent = 'Назначить';

  if (failed.length > 0) {
    errorEl.textContent = `Не удалось назначить студентам: ${failed.join(', ')}.`;
    errorEl.hidden = false;
  }

  const ok = selected.length - failed.length;
  if (ok > 0) {
    successEl.textContent = `Квиз назначен ${ok} студент${plural(ok, 'у', 'ам', 'ам')}.`;
    successEl.hidden = false;
    // uncheck assigned students
    document.querySelectorAll('.student-cb:checked').forEach(cb => {
      if (!failed.includes(parseInt(cb.value))) {
        cb.closest('.student-row').remove();
      }
    });
  }
}

function plural(n, one, few, many) {
  const mod10 = n % 10, mod100 = n % 100;
  if (mod10 === 1 && mod100 !== 11) return one;
  if (mod10 >= 2 && mod10 <= 4 && (mod100 < 10 || mod100 >= 20)) return few;
  return many;
}

function renderGeneratedQuiz(quiz) {
  const resultEl = document.getElementById('ai-result');
  resultEl.innerHTML = `
    <div class="generated-quiz">
      <div class="generated-quiz-header">
        <div>
          <h3>${esc(quiz.title)}</h3>
          <span class="quiz-meta">${quiz.questions?.length ?? 0} вопросов · ID ${quiz.id}</span>
        </div>
        <span class="status-badge status-done">Сохранён</span>
      </div>
      <div class="generated-questions">
        ${(quiz.questions ?? []).map((q, i) => `
          <div class="question-view">
            <p class="question-view-text">${i + 1}. ${esc(q.text)}</p>
            <ul class="options-view">
              ${(q.options ?? []).sort((a, b) => a.number - b.number).map(opt => `
                <li class="${opt.number === q.correct_answer_number ? 'opt-correct' : 'opt-normal'}">
                  ${opt.number}. ${esc(opt.text)}
                  ${opt.number === q.correct_answer_number ? '<span class="checkmark">✓</span>' : ''}
                </li>`).join('')}
            </ul>
          </div>`).join('')}
      </div>
    </div>`;
  resultEl.hidden = false;
  resultEl.scrollIntoView({ behavior: 'smooth', block: 'start' });
}

// ── localStorage ───────────────────────────────────────────
function getCompleted() {
  return JSON.parse(localStorage.getItem('completed') || '[]');
}

function markCompleted(id) {
  const list = getCompleted();
  if (!list.includes(id)) {
    list.push(id);
    localStorage.setItem('completed', JSON.stringify(list));
  }
}

function getSavedResults() {
  return JSON.parse(localStorage.getItem('results') || '[]');
}

function saveResult(assignmentId, quizTitle, attempt) {
  const results = getSavedResults();
  results.push({
    assignmentId,
    quizTitle,
    score:        attempt.score ?? 0,
    correctCount: attempt.correct_question_ids?.length ?? 0,
    wrongCount:   attempt.wrong_question_ids?.length   ?? 0,
    date:         new Date().toISOString(),
  });
  localStorage.setItem('results', JSON.stringify(results));
}

// ── Render helpers ─────────────────────────────────────────
function pageHeader(title, count) {
  return `<div class="page-header">
    <h2>${title}</h2>
    ${count != null ? `<span class="count-badge">${count}</span>` : ''}
  </div>`;
}

function errorState(msg) {
  return `<div class="error-state">${esc(msg)}</div>`;
}

function esc(str) {
  const d = document.createElement('div');
  d.textContent = str ?? '';
  return d.innerHTML;
}

// ── SVG icons ──────────────────────────────────────────────
function navIconQuizzes() {
  return `<svg width="18" height="18" viewBox="0 0 24 24" fill="none"
    stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
    <path d="M9 11l3 3L22 4"/><path d="M21 12v7a2 2 0 01-2 2H5a2 2 0 01-2-2V5a2 2 0 012-2h11"/>
  </svg>`;
}

function navIconCreate() {
  return `<svg width="18" height="18" viewBox="0 0 24 24" fill="none"
    stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
    <circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="16"/>
    <line x1="8" y1="12" x2="16" y2="12"/>
  </svg>`;
}

function navIconResults() {
  return `<svg width="18" height="18" viewBox="0 0 24 24" fill="none"
    stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
    <line x1="18" y1="20" x2="18" y2="10"/><line x1="12" y1="20" x2="12" y2="4"/>
    <line x1="6" y1="20" x2="6" y2="14"/>
  </svg>`;
}

// ── Start ──────────────────────────────────────────────────
init();
