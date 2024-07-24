// DOM elements
let searchInput;
let lifterTable;
let errorMessage;
let searchResults;

// Wait for the DOM to be fully loaded before accessing elements
document.addEventListener('DOMContentLoaded', function() {
    // Initialize variables after the DOM is loaded
    searchInput = document.getElementById('search-input');
    lifterTable = document.getElementById('lifter-table');
    errorMessage = document.getElementById('errorMessage');
    searchResults = document.getElementById('search-results');

    // If search-results doesn't exist, create and append it
    if (!searchResults) {
        searchResults = document.createElement('ul');
        searchResults.id = 'search-results';
        document.getElementById('search-container').appendChild(searchResults);
    }

    // Add event listener to the search input
    searchInput.addEventListener('input', debounce(searchLifters, 300));
});

function debounce(func, delay) {
    let debounceTimer;
    return function() {
        const context = this;
        const args = arguments;
        clearTimeout(debounceTimer);
        debounceTimer = setTimeout(() => func.apply(context, args), delay);
    }
}

function searchLifters() {
    const query = searchInput.value.trim();
    console.log("Search query:", query); // Add this line for debugging

    if (query.length === 0) {
        clearTable();
        searchResults.innerHTML = '';
        return;
    }

    fetch(`/api/search?q=${encodeURIComponent(query)}`)
        .then(response => response.json())
        .then(lifters => {
            console.log("Search results:", lifters); // Add this line for debugging
            searchResults.innerHTML = '';
            if (lifters.length > 0) {
                lifters.forEach(lifter => {
                    const li = document.createElement('li');
                    li.textContent = lifter;
                    li.addEventListener('click', () => displayLifterDetails(lifter));
                    searchResults.appendChild(li);
                });
            } else {
                errorMessage.textContent = 'No lifters found.';
            }
        })
        .catch(error => {
            console.error('Error:', error);
            errorMessage.textContent = 'An error occurred while searching.';
        });
}

function displayLifterDetails(lifterName) {
    fetch(`/api/lifter-details?name=${encodeURIComponent(lifterName)}`)
        .then(response => response.json())
        .then(details => {
            clearTable();
            const tbody = lifterTable.querySelector('tbody');
            details.forEach(detail => {
                const row = tbody.insertRow();
                row.innerHTML = `
                    <td>${detail.date}</td>
                    <td>${detail.meetName}</td>
                    <td>${detail.successfulSquatAttempts}</td>
                    <td>${detail.successfulBenchAttempts}</td>
                    <td>${detail.successfulDeadliftAttempts}</td>
                    <td>${detail.totalSuccessfulAttempts}</td>
                    <td>${detail.squat1Perc.toFixed(2)}</td>
                    <td>${detail.squat2Perc.toFixed(2)}</td>
                    <td>${detail.squat3Perc.toFixed(2)}</td>
                    <td>${detail.bench1Perc.toFixed(2)}</td>
                    <td>${detail.bench2Perc.toFixed(2)}</td>
                    <td>${detail.bench3Perc.toFixed(2)}</td>
                    <td>${detail.deadlift1Perc.toFixed(2)}</td>
                    <td>${detail.deadlift2Perc.toFixed(2)}</td>
                    <td>${detail.deadlift3Perc.toFixed(2)}</td>
                    <td>${detail.squat1To2Kg.toFixed(2)}</td>
                    <td>${detail.squat2To3Kg.toFixed(2)}</td>
                    <td>${detail.bench1To2Kg.toFixed(2)}</td>
                    <td>${detail.bench2To3Kg.toFixed(2)}</td>
                    <td>${detail.deadlift1To2Kg.toFixed(2)}</td>
                    <td>${detail.deadlift2To3Kg.toFixed(2)}</td>
                `;
            });
            searchResults.innerHTML = '';
        })
        .catch(error => {
            console.error('Error:', error);
            errorMessage.textContent = 'An error occurred while fetching lifter details.';
        });
}

function clearTable() {
    const tbody = lifterTable.querySelector('tbody');
    tbody.innerHTML = '';
    errorMessage.textContent = '';
}