import { createNamespacedHelpers } from 'vuex';
import { isFunction, isString } from 'lodash';

import popupMixin from '@/mixins/popup';

const { mapActions, mapGetters } = createNamespacedHelpers('webhook');

function mapActionsWith(actions, success, error) {
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
          return errorWithContext(err);
        }

        throw err;
      }
    };

    return acc;
  }, {});
}

function mapActionsWithPopup(actions, messages) {
  return mapActionsWith(actions, async function success(actionKey, result) {
    const text = messages[actionKey] && messages[actionKey].success ? messages[actionKey].success.apply(this) : this.$t('success.default');

    if (this.addSuccessPopup) {
      await this.addSuccessPopup({ text });
    }

    return result;
  }, async function error(popups) {
    const errorPopup = popups && popups.error ? popups.error : { text: this.$t('errors.default') };

    if (this.addErrorPopup) {
      await this.addErrorPopup(errorPopup);
    }
  });
}

function create(messages = {}) {
  return actions => mapActionsWithPopup(actions, messages);
}

const mapActionsWithPopupNew = create({
  createWebhookWithPopupWithRefresh: {
    success() {
      return 'Created!';
    },
    error() {
      return 'Create error!';
    },
  },
  updateWebhookWithPopupWithRefresh: {
    success() {
      return 'LOLOLO';
    },
    error() {
      return 'Update error!';
    },
  },
  removeWebhookWithPopupWithRefresh: {
    success() {
      return 'Removed!';
    },
    error() {
      return 'Remove error!';
    },
  },
});

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

    ...mapActionsWith(mapActionsWithPopupNew(mapActions({
      createWebhookWithPopupWithRefresh: 'create',
      updateWebhookWithPopupWithRefresh: 'update',
      removeWebhookWithPopupWithRefresh: 'remove',
    })), 'refreshWebhooksList'),
  },
};
