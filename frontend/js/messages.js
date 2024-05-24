let json;

async function bodyToJSON(data) {
    await json = data;
}

async function saveMessage() {
    const messageInput = document.getElementById('chatInput').value;

    fetch('/send-message', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(json)
    })
        .then(response => {
            if (!response.ok) {
                throw new Error('Failed to send message');
            }
            return response.json();
        })
        .then(data => {
            console.log('Message sent successfully:', data);
        })
        .catch(error => {
            console.error('Error sending message:', error.message);
        });
    document.getElementById('chatInput').value = '';
}
