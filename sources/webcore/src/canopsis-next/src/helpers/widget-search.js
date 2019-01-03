export function getContextSearchByText(text) {
  return {
    $or: [
      {
        name: { $regex: text, $options: 'i' },
      }, {
        type: { $regex: text, $options: 'i' },
      },
    ],
  };
}

export default {
  getContextSearchByText,
};
