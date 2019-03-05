import { isFunction, isString } from 'lodash';

export default function mapActionsWith(actions, success, error) {
  return Object.entries(actions).reduce((acc, [key, value]) => {
    acc[key] = async function mappedActionWithCalledAfter(...args) {
      try {
        let result;
        let successWithContext;

        if (success) {
          if (isFunction(success)) {
            successWithContext = success.bind(this);
          } else if (isString(success)) {
            successWithContext = this[success];
          }
        }

        if (isFunction(value)) {
          result = await value.apply(this, args);
        } else if (isString(value)) {
          result = await this[value](...args);
        }

        if (successWithContext) {
          return successWithContext(key, result);
        }

        return result;
      } catch (err) {
        let errorWithContext;

        if (error) {
          if (isFunction(error)) {
            errorWithContext = error.bind(this);
          } else if (isString(error)) {
            errorWithContext = this[error];
          }
        }

        if (errorWithContext) {
          return errorWithContext(key, err);
        }

        throw err;
      }
    };

    return acc;
  }, {});
}

export function createMapActionsWith(success, error) {
  return (actions, prioritySuccess, priorityError) =>
    mapActionsWith(actions, prioritySuccess || success, priorityError || error);
}
