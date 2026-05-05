const API_BASE = "/api/v1";

let currentUserId = 3;
let currentDate = new Date();
let currentView = "week";
let currentEventId = null;

document.addEventListener("DOMContentLoaded", () => {
  loadStateFromURL();
  setupListeners();
  loadEvents();
});

function saveStateToURL() {
  const params = new URLSearchParams();

  params.set("date", formatDate(currentDate));
  params.set("view", currentView);

  history.replaceState({}, "", `?${params.toString()}`);
}

function loadStateFromURL() {
  const params = new URLSearchParams(window.location.search);

  const date = params.get("date");
  const view = params.get("view");

  if (date) currentDate = new Date(date);
  if (view) currentView = view;
}

function formatDate(date) {
  const y = date.getFullYear();
  const m = String(date.getMonth() + 1).padStart(2, "0");
  const d = String(date.getDate()).padStart(2, "0");
  return `${y}-${m}-${d}`;
}

function buildRFC3339(dateStr, timeStr = "12:00") {
  const local = new Date(`${dateStr}T${timeStr}:00`);

  const pad = (n) => String(n).padStart(2, "0");

  const tzOffset = -local.getTimezoneOffset();
  const sign = tzOffset >= 0 ? "+" : "-";

  const hours = pad(Math.floor(Math.abs(tzOffset) / 60));
  const minutes = pad(Math.abs(tzOffset) % 60);

  const tz = `${sign}${hours}:${minutes}`;

  return `${dateStr}T${timeStr}:00${tz}`;
}

async function loadEvents() {
  const dateStr = formatDate(currentDate);

  let endpoint = "events_for_week";
  if (currentView === "day") endpoint = "events_for_day";
  if (currentView === "month") endpoint = "events_for_month";

  const url = `${API_BASE}/${endpoint}?user_id=${currentUserId}&date=${dateStr}`;

  try {
    const res = await fetch(url);
    const data = await res.json();

    if (!res.ok) {
      alert(data.error || "Loading error");
      return;
    }

    renderEvents(data.result?.events || []);
  } catch (err) {
    console.error(err);
    alert("Server unavailable");
  }
}

function renderEvents(events) {
  const container = document.getElementById("events-list");
  container.innerHTML = "<h3>Events</h3>";

  if (!events.length) {
    container.innerHTML += "<p>No events</p>";
    return;
  }

  events.forEach((ev) => {
    const div = document.createElement("div");
    div.className = "event-item";

    div.innerHTML = `
            <strong>${ev.date}</strong><br>
            ${ev.text}
        `;

    div.onclick = () => openModal(ev);

    container.appendChild(div);
  });
}

function setupListeners() {
  document.getElementById("btn-new-event").onclick = () => openModal();
  document.getElementById("btn-cancel").onclick = closeModal;
  document.getElementById("event-form").onsubmit = saveEvent;
  document.getElementById("btn-delete").onclick = deleteEvent;

  document.getElementById("btn-today").onclick = () => {
    currentDate = new Date();
    saveStateToURL();
    loadEvents();
  };

  document.getElementById("btn-prev").onclick = () => navigate(-1);
  document.getElementById("btn-next").onclick = () => navigate(1);

  document.getElementById("view-mode").value = currentView;
  document.getElementById("view-mode").onchange = (e) => {
    currentView = e.target.value;
    saveStateToURL();
    loadEvents();
  };
}

function openModal(event = null) {
  currentEventId = null;

  const modal = document.getElementById("event-modal");
  const deleteBtn = document.getElementById("btn-delete");

  if (event) {
    currentEventId = event.event_id;

    document.getElementById("event-date").value = event.date.slice(0, 10);
    document.getElementById("event-text").value = event.text;

    deleteBtn.style.display = "block";
  } else {
    document.getElementById("event-form").reset();
    document.getElementById("event-date").value = formatDate(currentDate);

    deleteBtn.style.display = "none";
  }

  modal.style.display = "flex";
}

function closeModal() {
  document.getElementById("event-modal").style.display = "none";
}

async function saveEvent(e) {
  e.preventDefault();

  const date = document.getElementById("event-date").value;
  const time = document.getElementById("event-time").value || "12:00";
  const text = document.getElementById("event-text").value;
  const reminder = Number(document.getElementById("event-reminder").value || 0);

  const fullDate = buildRFC3339(date, time);

  let url = `${API_BASE}/create_event`;
  let payload = {
    user_id: currentUserId,
    date: fullDate,
    text,
    reminder,
  };

  if (currentEventId) {
    url = `${API_BASE}/update_event`;
    payload = {
      user_id: currentUserId,
      event_id: currentEventId,
      text,
      new_date: fullDate,
    };
  }

  try {
    const res = await fetch(url, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });

    const data = await res.json();

    if (!res.ok) {
      alert(data.error || "Save error");
      return;
    }

    closeModal();
    saveStateToURL();
    loadEvents();
  } catch (err) {
    console.error(err);
    alert("Network error");
  }
}

async function deleteEvent() {
  if (!currentEventId) return;
  if (!confirm("Delete event?")) return;

  try {
    const res = await fetch(`${API_BASE}/delete_event`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        user_id: currentUserId,
        event_id: currentEventId,
      }),
    });

    const data = await res.json();

    if (!res.ok) {
      alert(data.error || "Delete error");
      return;
    }

    closeModal();
    loadEvents();
  } catch (err) {
    console.error(err);
    alert("Network error");
  }
}

function navigate(dir) {
  if (currentView === "day") {
    currentDate.setDate(currentDate.getDate() + dir);
  } else if (currentView === "week") {
    currentDate.setDate(currentDate.getDate() + 7 * dir);
  } else {
    currentDate.setMonth(currentDate.getMonth() + dir);
  }

  saveStateToURL();
  loadEvents();
}
