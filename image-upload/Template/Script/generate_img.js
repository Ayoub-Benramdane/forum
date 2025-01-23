document.addEventListener("DOMContentLoaded", function () {
    const form = document.querySelector("#new-post-form");

    form.addEventListener("submit", async function (event) {
        event.preventDefault();

        const formData = new FormData(form);

        try {
            const response = await fetch(form.action, {
                method: "POST",
                body: formData
            });

            if (response.ok) {
                window.location.href = "/";
            } else if (response.status === 406) {
                alert("Only image files are allowed or the file is too large.");
            } else {
                alert("Failed to post the file.");
            }
        } catch (error) {
            alert("An error occurred while uploading the file. Please try again.");
            console.error("Error:", error);
        }
    });
});