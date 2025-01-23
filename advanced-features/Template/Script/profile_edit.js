document.addEventListener("DOMContentLoaded", function () {
    const form = document.querySelector("form");
    const usernameInput = document.getElementById("username");
    const emailInput = document.getElementById("email");
    const passwordInput = document.getElementById("password");
    const newPasswordInput = document.getElementById("new-password");
    const confirmPasswordInput = document.getElementById("confirm-password");
    const generalError = document.getElementById("generalError");

    form.addEventListener("submit", async function (e) {
        e.preventDefault();
        clearErrors();

        let isValid = true;

        if (usernameInput.value.trim() === "") {
            showError(usernameInput, "Username is required.");
            isValid = false;
        }

        if (!validateEmail(emailInput.value)) {
            showError(emailInput, "Please enter a valid email address.");
            isValid = false;
        }

        if (newPasswordInput.value !== "") {
            if (!validatePassword(newPasswordInput.value)) {
                showError(newPasswordInput, "Password must be at least 8 characters, include an uppercase letter, a number, and a special character.");
                isValid = false;
            }

            if (newPasswordInput.value !== confirmPasswordInput.value) {
                showError(confirmPasswordInput, "Passwords do not match.");
                isValid = false;
            }
        }

        if (!isValid) {
            return;
        }

        try {
            const response = await fetch("/profile_edit", {
                method: "POST",
                headers: {
                    "Content-Type": "application/x-www-form-urlencoded",
                },
                body: `password=${encodeURIComponent(passwordInput.value)}`,
            });

            if (response.ok) {
                window.location.href = "/profile";
            } else if (response.status === 401) {
                generalError.textContent = "Invalid password.";
            } else {
                generalError.textContent = "An unexpected error occurred. Please try again.";
            }
        } catch (error) {
            generalError.textContent = "Unable to connect to the server. Please try again.";
        }
    });

    function validateEmail(email) {
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        return emailRegex.test(email);
    }

    function validatePassword(password) {
        const passwordRegex = /^(?=.*[A-Z])(?=.*\d)(?=.*[!@#$%^&*])[A-Za-z\d!@#$%^&*]{8,}$/;
        return passwordRegex.test(password);
    }

    function showError(input, message) {
        const errorElement = document.createElement("span");
        errorElement.className = "error-message";
        errorElement.textContent = message;
        input.parentElement.appendChild(errorElement);
        input.classList.add("error");
    }

    function clearErrors() {
        const errorMessages = document.querySelectorAll(".error-message");
        errorMessages.forEach((msg) => msg.remove());

        const inputs = document.querySelectorAll(".error");
        inputs.forEach((input) => input.classList.remove("error"));

        generalError.textContent = "";
    }
});