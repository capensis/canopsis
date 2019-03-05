import { isFunction, isString } from 'lodash';

export default function (actions, success, error) {
  return Object.entries(actions).reduce((acc, [key, value]) => {
    acc[key] = async function mappedActionWithCalledAfter(...args) {
      try {
        let result;
        let actionArgs = args;
        let successArgs = [];
        let successWithContext;

        if (success) {
          if (isFunction(success)) {
            successWithContext = success.bind(this);
          } else if (isString(success)) {
            successWithContext = this[success];
          }
        }

        if (successWithContext && successWithContext.length) {
          actionArgs = args.slice(successWithContext.length);
          successArgs = args.slice(0, successWithContext.length);
        }


        if (isFunction(value)) {
          result = await value.apply(this, actionArgs);
        } else if (isString(value)) {
          result = await this[value](...actionArgs);
        }

        if (successWithContext) {
          await successWithContext(...successArgs);
        }

        return result;
      } catch (err) {
        let errorArgs = [];
        let errorWithContext;

        if (error) {
          if (isFunction(error)) {
            errorWithContext = error.bind(this);
          } else if (isString(error)) {
            errorWithContext = this[error];
          }
        }

        if (errorWithContext && errorWithContext.length) {
          errorArgs = args.slice(0, errorWithContext.length);
        }

        if (errorWithContext) {
          await errorWithContext(...errorArgs);
        }

        throw err;
      }
    };

    return acc;
  }, {});
}
