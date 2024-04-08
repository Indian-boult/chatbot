document.addEventListener('DOMContentLoaded', function() {
    const form = document.getElementById('upload-form');

    form.addEventListener('submit', async function(e) {
        e.preventDefault();
        const formData = new FormData(this);

        try {
            const response = await fetch('http://localhost:8080/upload', {
                method: 'POST',
                body: formData,
            });

            const messageBox = document.getElementById('chat-box');
            let messageText;

            if (response.ok) {
                const data = await response.text(); // Assuming the server sends back a plain text response
                messageText = 'Image uploaded successfully: ' + data; // data contains the server response
            } else {
                const error = await response.text();
                messageText = 'Failed to upload image: ' + error;
            }

            const message = document.createElement('div');
            message.textContent = messageText;
            messageBox.appendChild(message);
        } catch (error) {
            console.error('Upload failed:', error);
        }

        form.reset(); // Resets the form after the upload attempt
    });
});
