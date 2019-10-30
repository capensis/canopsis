export default (store) => {
  // eslint-disable-next-line no-param-reassign
  store.watchOnce = (getter, comparator = v => v) => new Promise((resolve) => {
    const unwatch = store.watch(getter, (value) => {
      if (comparator(value)) {
        unwatch();

        resolve(value);
      }
    });
  });
};
