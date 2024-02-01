document.getElementById('accountDropdownBtn').addEventListener('mouseenter', function () {
  document.getElementById('accountDropdownMenu').classList.add('show');
});

document.getElementById('accountDropdownBtn').addEventListener('mouseleave', function () {
  document.getElementById('accountDropdownMenu').classList.remove('show');
});

document.getElementById('accountDropdownMenu').addEventListener('mouseenter', function () {
  this.classList.add('show');
});

document.getElementById('accountDropdownMenu').addEventListener('mouseleave', function () {
  this.classList.remove('show');
});
