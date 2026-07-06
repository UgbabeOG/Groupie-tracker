// ============================================================
// Groupie Tracker
// Main JavaScript File
//
// Features:
// - Search artists
// - Sort artists
// - Fade-in animation
// - Hover effects
// - Load More button
// ============================================================

document.addEventListener("DOMContentLoaded", () => {

    // ========================================================
    // ELEMENTS
    // ========================================================

    const searchInput = document.getElementById("artist-search");
    const sortSelect = document.getElementById("sort-artists");
    const loadMoreBtn = document.getElementById("load-more");
    const grid = document.getElementById("artist-grid");

    // Get all artist cards
    const cards = Array.from(document.querySelectorAll(".artist-card"));

    // ========================================================
    // INITIAL CARD ANIMATION
    // ========================================================

    cards.forEach((card, index) => {

        card.style.opacity = "0";
        card.style.transform = "translateY(20px)";

        setTimeout(() => {

            card.style.opacity = "1";
            card.style.transform = "translateY(0)";

        }, index * 80);

    });

    // ========================================================
    // SEARCH
    // ========================================================

    if (searchInput) {

        searchInput.addEventListener("input", () => {

            const search = searchInput.value.trim().toLowerCase();

            cards.forEach(card => {

                const name = card
                    .querySelector(".card-title")
                    .textContent
                    .toLowerCase();

                if (name.includes(search)) {

                    card.style.display = "flex";

                } else {

                    card.style.display = "none";

                }

            });

        });

    }

    // ========================================================
    // HOVER EFFECT
    // ========================================================

    cards.forEach(card => {

        card.addEventListener("mouseenter", () => {

            card.style.transform = "translateY(-8px)";

        });

        card.addEventListener("mouseleave", () => {

            card.style.transform = "translateY(0)";

        });

    });

    // ========================================================
    // LOAD MORE
    // ========================================================

    if (loadMoreBtn) {

        let visibleCards = 8;

        // Hide cards after the first 8
        cards.forEach((card, index) => {

            if (index >= visibleCards) {

                card.style.display = "none";

            }

        });

        loadMoreBtn.addEventListener("click", () => {

            const text = loadMoreBtn.querySelector(".btn-text");

            if (text) {

                text.textContent = "Loading...";

            }

            setTimeout(() => {

                visibleCards += 8;

                cards.forEach((card, index) => {

                    if (index < visibleCards) {

                        card.style.display = "flex";

                    }

                });

                if (text) {

                    text.textContent = "Load More";

                }

                // Hide button when all artists are visible
                if (visibleCards >= cards.length) {

                    loadMoreBtn.style.display = "none";

                }

            }, 400);

        });

    }

    // ========================================================
    // SORT
    // ========================================================

    if (sortSelect && grid) {

        sort.addEventListener("change", () => {

            const cards = Array.from(document.querySelectorAll(".artist-card"));

            switch (sort.value) {

                case "name":

                    cards.sort((a, b) =>
                        a.dataset.name.localeCompare(b.dataset.name)
                    );
                    break;

                case "members":

                    cards.sort((a, b) =>
                        Number(b.dataset.members) - Number(a.dataset.members)
                    );
                    break;

                case "year":

                    cards.sort((a, b) =>
                        Number(a.dataset.year) - Number(b.dataset.year)
                    );
                    break;
            }

            cards.forEach(card => grid.appendChild(card));

        });

    }

});