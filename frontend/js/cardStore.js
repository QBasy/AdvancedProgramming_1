const cardsPerPage = 12;
let currentPage = 1;

function createPageSelector(totalPages) {
  const pageSelector = document.createElement('div');
  pageSelector.classList.add('page-selector');
  for (let i = 1; i <= totalPages; i++) {
    const pageButton = document.createElement('button');
    pageButton.textContent = i;
    pageButton.addEventListener('click', () => {
      showPage(i);
    });
    pageSelector.appendChild(pageButton);
  }
  return pageSelector;
}

function showPage(pageNumber) {
  const start = (pageNumber - 1) * cardsPerPage;
  const end = pageNumber * cardsPerPage;
  const pageCards = jsonData.slice(start, end);

  const section = document.querySelector('.py-5 .container .row');
  section.innerHTML = '';
  pageCards.forEach(cardData => {
    const card = createCard(cardData.title, cardData.author, cardData.date, cardData.imageURL);
    section.appendChild(card);
  });

  const pageSelector = document.querySelector('.page-selector');
  if (pageSelector) {
    pageSelector.remove();
  }
  if (Math.ceil(jsonData.length / cardsPerPage) > 1) {
    section.parentElement.appendChild(createPageSelector(Math.ceil(jsonData.length / cardsPerPage)));
  }
}

function createCard(title, author, date, imageURL) {
  const col = document.createElement('div');
  col.classList.add('col', 'mb-5');

  const card = document.createElement('div');
  card.classList.add('card', 'h-100');

  const img = document.createElement('img');
  img.classList.add('card-img-top');
  img.src = imageURL;
  img.alt = title;

  const cardBody = document.createElement('div');
  cardBody.classList.add('card-body', 'p-4');

  const cardTitle = document.createElement('h5');
  cardTitle.classList.add('fw-bolder');
  cardTitle.textContent = title;

  const cardSubtitle = document.createElement('h6');
  cardSubtitle.textContent = author;

  const cardDate = document.createElement('p');
  cardDate.classList.add('text-end', 'mb-0');
  cardDate.textContent = date;

  cardBody.appendChild(cardTitle);
  cardBody.appendChild(cardSubtitle);
  cardBody.appendChild(cardDate);

  card.appendChild(img);
  card.appendChild(cardBody);

  col.appendChild(card);

  return col;
}

const jsonData = [
  {
    "title": "Jujutsu Kaisen S1EP2",
    "author": "megaTerminator228",
    "date": "2 weeks ago",
    "imageURL": "https://dummyimage.com/450x300/dee2e6/6c757d.jpg"
  },
  {
    "title": "Interstellar",
    "author": "moviePirate1337",
    "date": "5 days ago",
    "imageURL": "https://dummyimage.com/450x300/dee2e6/6c757d.jpg"
  },
  {
    "title": "Markiplier on FlickFusion???",
    "author": "Markiplier",
    "date": "3 weeks ago",
    "imageURL": "https://dummyimage.com/450x300/dee2e6/6c757d.jpg"
  },
  {
    "title": "TOP 50 ANIME",
    "author": "ANIMELOVER",
    "date": "1 day ago",
    "imageURL": "https://dummyimage.com/450x300/dee2e6/6c757d.jpg"
  },
  {
    "title": "New Card 1",
    "author": "Author 1",
    "date": "2 days ago",
    "imageURL": "https://dummyimage.com/450x300/dee2e6/6c757d.jpg"
  },
  {
    "title": "New Card 2",
    "author": "Author 2",
    "date": "3 days ago",
    "imageURL": "https://dummyimage.com/450x300/dee2e6/6c757d.jpg"
  },
  {
    "title": "New Card 3",
    "author": "Author 3",
    "date": "3 days ago",
    "imageURL": "https://dummyimage.com/450x300/dee2e6/6c757d.jpg"
  },
  {
    "title": "New Card 4",
    "author": "Author 4",
    "date": "1 day ago",
    "imageURL": "https://dummyimage.com/450x300/dee2e6/6c757d.jpg"
  },
  {
    "title": "New Card 5",
    "author": "Author 5",
    "date": "2 days ago",
    "imageURL": "https://dummyimage.com/450x300/dee2e6/6c757d.jpg"
  },
  {
    "title": "New Card 6",
    "author": "Author 6",
    "date": "3 days ago",
    "imageURL": "https://dummyimage.com/450x300/dee2e6/6c757d.jpg"
  },
  {
    "title": "New Card 7",
    "author": "Author 7",
    "date": "3 days ago",
    "imageURL": "https://dummyimage.com/450x300/dee2e6/6c757d.jpg"
  },
  {
    "title": "New Card 8",
    "author": "Author 8",
    "date": "1 day ago",
    "imageURL": "https://dummyimage.com/450x300/dee2e6/6c757d.jpg"
  },
  {
    "title": "New Card 9",
    "author": "Author 1",
    "date": "2 days ago",
    "imageURL": "https://dummyimage.com/450x300/dee2e6/6c757d.jpg"
  },
  {
    "title": "New Card 10",
    "author": "Author 2",
    "date": "3 days ago",
    "imageURL": "https://dummyimage.com/450x300/dee2e6/6c757d.jpg"
  },
  {
    "title": "New Card 11",
    "author": "Author 3",
    "date": "3 days ago",
    "imageURL": "https://dummyimage.com/450x300/dee2e6/6c757d.jpg"
  },
  {
    "title": "New Card 12",
    "author": "Author 4",
    "date": "1 day ago",
    "imageURL": "https://dummyimage.com/450x300/dee2e6/6c757d.jpg"
  },
];

showPage(currentPage);
