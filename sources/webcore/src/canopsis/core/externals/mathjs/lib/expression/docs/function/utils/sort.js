module.exports = {
  'name': 'sort',
  'category': 'Utils',
  'syntax': [
    'sort(x)',
    'sort(x, compare)'
  ],
  'description': 'Sort the items in a matrix. Compare can be a string "asc" or "desc", or a custom sort function.',
  'examples': [
    'sort([5, 10, 1])',
    'sort(["C", "B", "A", "D"])',
    'sortByLength(a, b) = size(a)[1] - size(b)[1]',
    'sort(["Langdon", "Tom", "Sara"], sortByLength)'
  ],
  'seealso': ['map', 'filter', 'forEach']
};
