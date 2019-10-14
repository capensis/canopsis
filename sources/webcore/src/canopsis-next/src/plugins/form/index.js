import { isEqual } from 'lodash';

import { setIn, unsetIn } from '@/helpers/immutable';

export default {
  install(Vue) {
    Object.defineProperty(Vue.prototype, '$form', {
      get() {
        const self = this;

        return {
          /**
           * Getter for model property name
           *
           * @returns {string}
           */
          get prop() {
            const { prop = 'value' } = self.$options.model || {};

            return prop;
          },

          /**
           * Getter for model event name
           *
           * @returns {string}
           */
          get event() {
            const { event = 'input' } = self.$options.model || {};

            return event;
          },

          /**
           * Update full model
           *
           * @param {*} value
           */
          updateModel: (value) => {
            this.$emit(this.$form.event, value);
          },

          /**
           * Emit event to parent with new object and with updated field
           *
           * @param {string|Array} path
           * @param {*} value
           */
          updateField: (path, value) => {
            this.$form.updateModel(path.length ? setIn(this[this.$form.prop], path, value) : value);
          },

          /**
           * Emit event to parent with new object
           * Rename a field in the object and update it
           *
           * @param {string|Array} oldPath
           * @param {string|Array} newPath
           * @param {*} value
           */
          updateAndMoveField(oldPath, newPath, value) {
            if (isEqual(oldPath, newPath)) {
              this.$form.updateField(oldPath, value);
            } else {
              const result = unsetIn(this[this.$form.prop], oldPath);

              this.$form.moveField(setIn(result, newPath, value));
            }
          },

          /**
           * Emit event to parent with new object
           * Rename a field in the object
           *
           * @param {string|Array} oldPath
           * @param {string|Array} newPath
           */
          moveField(oldPath, newPath) {
            if (!isEqual(oldPath, newPath)) {
              const result = unsetIn(this[this.$form.prop], oldPath);

              this.$form.updateModel(setIn(result, newPath, value => value));
            }
          },

          /**
           * Emit event to parent with new object without field
           *
           * @param {string|Array} path
           */
          removeField: (path) => {
            this.$form.updateModel(unsetIn(this[this.$form.prop], path));
          },

          /**
           * Emit event to parent with new array with new item
           *
           * @param {*} value
           */
          addItemIntoArray: (value) => {
            this.$form.updateModel([...this[this.$form.prop], value]);
          },

          /**
           * Emit event to parent with new array with updated array item
           *
           * @param {number} index
           * @param {*} value
           */
          updateItemInArray: (index, value) => {
            const items = [...this[this.$form.prop]];

            items[index] = value;

            this.$form.updateModel(items);
          },

          /**
           * Emit event to parent with new array with updated array item with updated field
           *
           * @param {number} index
           * @param {string} fieldName
           * @param {*} value
           */
          updateFieldInArrayItem: (index, fieldName, value) => {
            this.$form.updateItemInArray(index, setIn(this[this.$form.prop][index], fieldName, value));
          },

          /**
           * Emit event to parent with new array without array item
           *
           * @param {number} index
           */
          removeItemFromArray: (index) => {
            this.$form.updateModel(this[this.$form.prop].filter((v, i) => i !== index));
          },
        };
      },
    });
  },
};
