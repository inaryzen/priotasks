// Function to show the modal
function showModal(id) {
    document.getElementById(id).style.display = 'flex';
}

// Function to close the modal
function closeModal(id) {
    document.getElementById(id).style.display = 'none';
}

// Close modal when clicking outside the modal content
window.onclick = function (event) {
    const modal = document.getElementById('modal');
    if (event.target === modal) {
        closeModal();
    }
}