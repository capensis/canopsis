module.exports = {
  'name': 'bignumber',
  'category': 'Type',
  'syntax': [
    'bignumber(x)'
  ],
  'description':
      'Create a big number from a number or string.',
  'examples': [
    '0.1 + 0.2',
    'bignumber(0.1) + bignumber(0.2)',
    'bignumber("7.2")',
    'bignumber("7.2e500")',
    'bignumber([0.1, 0.2, 0.3])'
  ],
  'seealso': [
    'boolean', 'complex', 'index', 'matrix', 'string', 'unit'
  ]
};
