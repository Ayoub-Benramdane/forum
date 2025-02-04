document.addEventListener("DOMContentLoaded", function () {
    const form = document.querySelector("form");
    const titleInput = document.getElementById("title");
    const categorySelect = document.getElementById("category");
    const contentTextarea = document.getElementById("content");

    form.addEventListener("submit", function (event) {
        let isValid = true;
        clearErrors();

        if (titleInput && titleInput.value.trim() === "") {
            showError(titleInput, "Title is required.");
            isValid = false;
        }

        if (categorySelect && categorySelect.selectedOptions.length === 0) {
            showError(categorySelect, "Please select at least one category.");
            isValid = false;
        }

        if (contentTextarea.value.trim() === "") {
            showError(contentTextarea, "Content cannot be empty.");
            isValid = false;
        } else if (contentTextarea.value.trim().length < 10) {
            showError(contentTextarea, "Content must be at least 10 characters long.");
            isValid = false;
        }

        if (!isValid) {
            event.preventDefault();
        }
    });

    function showError(element, message) {
        const errorDiv = document.createElement("div");
        errorDiv.className = "error-message";
        errorDiv.innerText = message;
        element.classList.add("input-error");
        element.parentElement.appendChild(errorDiv);
    }

    function clearErrors() {
        const errors = document.querySelectorAll(".error-message");
        errors.forEach(error => error.remove());
        const inputs = document.querySelectorAll(".input-error");
        inputs.forEach(input => input.classList.remove("input-error"));
    }
});
