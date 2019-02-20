import uid from './uid';

export function textPairsToObject(textPairs) {
  return textPairs.reduce((acc, { text, value }) => {
    acc[text] = value;

    return acc;
  }, {});
}

export function objectToTextPairs(object) {
  return Object.keys(object).map(text => ({ key: uid(), text, value: object[text] }));
}

export function defaultItemCreator() {
  return { key: uid(), text: '', value: '' };
}

export default {
  textPairsToObject,
  objectToTextPairs,
  defaultItemCreator,
};
