document.getElementById("registerForm").addEventListener("submit", async function (e) {
    e.preventDefault();

    const errors = document.querySelectorAll(".error-message");
    errors.forEach(error => error.textContent = "");

    const username = document.getElementById("username").value.trim();
    const email = document.getElementById("email").value.trim();
    const password = document.getElementById("password").value.trim();
    const confirmPassword = document.getElementById("confirmPassword").value.trim();

    let hasError = false;

    if (username.length < 3 || username.length > 20) {
        document.getElementById("usernameError").textContent = "Username must be between 3 and 20 characters.";
        hasError = true;
    } else if (!/^[a-zA-Z0-9_]+$/.test(username)) {
        document.getElementById("usernameError").textContent = "Username can only contain letters, numbers, and underscores.";
        hasError = true;
    }

    if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
        document.getElementById("emailError").textContent = "Please enter a valid email address.";
        hasError = true;
    }

    if (password.length < 8) {
        document.getElementById("passwordError").textContent = "Password must be at least 8 characters long.";
        hasError = true;
    } else if (!/[A-Z]/.test(password) || !/[a-z]/.test(password) || !/[0-9]/.test(password) || !/[^a-zA-Z0-9]/.test(password)) {
        document.getElementById("passwordError").textContent = "Password must include uppercase, lowercase, number, and special character.";
        hasError = true;
    }

    if (password !== confirmPassword) {
        document.getElementById("confirmPasswordError").textContent = "Passwords do not match.";
        hasError = true;
    }

    if (hasError) return;

    try {
        const response = await fetch("/register", {
            method: "POST",
            headers: {
                "Content-Type": "application/x-www-form-urlencoded",
            },
            body: `username=${encodeURIComponent(username)}&email=${encodeURIComponent(email)}&password=${encodeURIComponent(password)}&confirm-password=${encodeURIComponent(confirmPassword)}`,
        });

        if (response.ok) {
            window.location.href = "/login";
        } else if (response.status === 409) {
            document.getElementById("generalError").textContent = "Username or email is already taken.";
        } else if (response.status === 400) {
            document.getElementById("generalError").textContent = "Invalid input. Please correct errors.";
        } else {
            document.getElementById("generalError").textContent = "An unexpected error occurred. Please try again.";
        }
    } catch (error) {
        document.getElementById("generalError").textContent = "Unable to connect to the server. Please try again.";
    }
});