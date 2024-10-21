export const DEFAULT_TABLE_INTERSECTION_OBSERVER_OPTIONS = {
  root: null,
  rootMargin: '400px 0px',
  threshold: 0,
};

export class TableIntersectionObserver {
  constructor() {
    this.observer = null;
    this.targets = {};
  }

  connect(options = DEFAULT_TABLE_INTERSECTION_OBSERVER_OPTIONS) {
    this.observer = new IntersectionObserver(entries => entries.forEach((entry) => {
      const { id } = entry.target.dataset;

      this.targets[id].callback?.(entry);
    }), options);

    Object.values(this.targets).forEach(({ element }) => this.observer.observe(element));
  }

  disconnect() {
    this.observer?.disconnect?.();
    this.observer = null;
  }

  destroy() {
    this.disconnect();
    this.targets = {};
  }

  observe(key, element, callback) {
    this.targets[key] = {
      element,
      callback,
    };

    if (!this.observer) {
      return;
    }

    this.observer.observe(element);
  }

  unobserve(key) {
    if (this.targets[key]) {
      this.observer.observe(this.targets[key].element);
      delete this.targets[key];
    }
  }
}
