import { createNamespacedHelpers } from 'vuex';
import { isFunction, isString } from 'lodash';

import popupMixin from '@/mixins/popup';

const { mapActions, mapGetters } = createNamespacedHelpers('webhook');

function mapActionsWith(actions, success, error) {
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

function mapActionsWithPopup(actions) {
  return mapActionsWith(actions, async function success(popups) {
    let successPopup;

    if (popups) {
      successPopup = popups.success ? popups.success : popups;
    } else {
      successPopup = { text: this.$t('success.default') };
    }

    if (this.addSuccessPopup) {
      await this.addSuccessPopup(successPopup);
    }
  }, async function error(popups) {
    const errorPopup = popups && popups.error ? popups.error : { text: this.$t('errors.default') };

    if (this.addErrorPopup) {
      await this.addErrorPopup(errorPopup);
    }
  });
}

export default {
  mixins: [popupMixin],
  computed: {
    ...mapGetters({
      webhooksPending: 'pending',
      webhooks: 'items',
    }),
  },
  methods: {
    ...mapActions({
      fetchWebhooksList: 'fetchList',
      refreshWebhooksList: 'fetchListWithPreviousParams',
    }),

    ...mapActionsWith(mapActionsWithPopup(mapActions({
      createWebhookWithPopupWithRefresh: 'create',
      updateWebhookWithPopupWithRefresh: 'update',
      removeWebhookWithPopupWithRefresh: 'remove',
    })), 'refreshWebhooksList'),
  },
};
