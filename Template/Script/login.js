document.getElementById("loginForm").addEventListener("submit", async function (e) {
    e.preventDefault();

    document.getElementById("usernameError").textContent = "";
    document.getElementById("passwordError").textContent = "";
    document.getElementById("generalError").textContent = "";

    const username = document.getElementById("username").value.trim();
    const password = document.getElementById("password").value.trim();

    let hasError = false;
    if (username.length < 3) {
        document.getElementById("usernameError").textContent = "Username must be at least 3 characters.";
        hasError = true;
    }
    if (password.length < 8) {
        document.getElementById("passwordError").textContent = "Password must be at least 8 characters.";
        hasError = true;
    }

    if (hasError) return;

    try {
        const response = await fetch("/login", {
            method: "POST",
            headers: {
                "Content-Type": "application/x-www-form-urlencoded",
            },
            body: `username=${encodeURIComponent(username)}&password=${encodeURIComponent(password)}`,
        });

        if (response.ok) {
            window.location.href = "/";
        } else if (response.status === 401) {
            document.getElementById("generalError").textContent = "Invalid username or password.";
        } else {
            document.getElementById("generalError").textContent = "An unexpected error occurred. Please try again.";
        }
    } catch (error) {
        document.getElementById("generalError").textContent = "Unable to connect to the server. Please try again.";
    }
});