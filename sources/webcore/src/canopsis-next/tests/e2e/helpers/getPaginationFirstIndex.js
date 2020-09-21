/**
 * Get first element index for pagination
 *
 * @param {Number} page
 * @param {Number} perPage
 * @returns {number}
 */
module.exports = (page, perPage) => ((page - 1) * perPage) + 1;
