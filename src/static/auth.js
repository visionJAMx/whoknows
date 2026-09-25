document.addEventListener("DOMContentLoaded", function () {
  var forms = document.querySelectorAll("#login-form, #register-form");

  forms.forEach(function (form) {
    form.addEventListener("submit", function (event) {
      event.preventDefault();
      handleSubmit(form);
    });
  });
});

function handleSubmit(form) {
  var errorEl = form.querySelector('[data-role="form-error"]');
  var formData = new FormData(form);
  var params = new URLSearchParams();
  formData.forEach(function (value, key) {
    params.append(key, value);
  });

  hideError(errorEl);

  fetch(form.getAttribute("action"), {
    method: "POST",
    headers: {
      "Content-Type": "application/x-www-form-urlencoded",
    },
    body: params.toString(),
  })
    .then(function (response) {
      return response.json().then(function (data) {
        return { ok: response.ok, data: data };
      });
    })
    .then(function (result) {
      if (result.ok) {
        redirectAfterSuccess(form);
      } else {
        showError(errorEl, extractErrorMessage(result.data));
      }
    })
    .catch(function () {
      showError(errorEl, "Der opstod en uventet fejl. Prøv igen.");
    });
}

function extractErrorMessage(data) {
  if (data && Array.isArray(data.detail) && data.detail.length > 0) {
    return data.detail
      .map(function (item) {
        return item.msg;
      })
      .join(" ");
  }
  if (data && data.message) {
    return data.message;
  }
  return "Der opstod en fejl. Prøv igen.";
}

function showError(errorEl, message) {
  if (!errorEl) return;
  errorEl.textContent = message;
  errorEl.hidden = false;
}

function hideError(errorEl) {
  if (!errorEl) return;
  errorEl.textContent = "";
  errorEl.hidden = true;
}

function redirectAfterSuccess(form) {
  if (form.id === "login-form") {
    window.location.href = "/";
  } else if (form.id === "register-form") {
    window.location.href = "/login";
  }
}
