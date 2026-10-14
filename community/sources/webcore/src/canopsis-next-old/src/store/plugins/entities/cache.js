/* eslint-disable no-underscore-dangle */
import { get } from 'lodash';

class Cache {
  constructor(init) {
    this._wm = new WeakMap(init);
  }

  clear() {
    this._wm = new WeakMap();
  }

  delete(k) {
    return this._wm.delete(k);
  }

  get(k) {
    return this._wm.get(k);
  }

  has(k) {
    return this._wm.has(k);
  }

  set(k, v) {
    this._wm.set(k, v);

    return this;
  }

  clearForEntityParents(state = {}, parents = {}) {
    parents.forEach((parent) => {
      const entity = get(state, [parent.type, parent.id]);

      if (entity) {
        this.delete(entity);

        if (entity._embedded && entity._embedded.parents) {
          this.clearForEntityParents(state, entity._embedded.parents);
        }
      }
    });
  }

  clearForEntity(state = {}, entity) {
    this.delete(entity);

    if (entity._embedded && entity._embedded.parents) {
      this.clearForEntityParents(state, entity._embedded.parents);
    }
  }
}

export default new Cache();
