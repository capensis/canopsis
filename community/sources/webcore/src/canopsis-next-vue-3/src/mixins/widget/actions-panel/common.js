import { isEmpty, pick } from 'lodash';
import { mapActions } from 'pinia';

import { PBEHAVIOR_ORIGINS } from '@/constants';

import { usePbehaviorStore } from '@/stores/modules/pbehavior';

import { createDowntimePbehavior } from '@/helpers/entities/pbehavior/form';

export const widgetActionsPanelCommonMixin = {
  methods: {
    ...mapActions(usePbehaviorStore, ['bulkCreateEntityPbehaviors', 'bulkRemoveEntityPbehaviors']),

    showPbehaviorResponsePopups(response, remove = false) {
      if (response?.length) {
        let hasErrors;

        response.forEach(({ error, errors }) => {
          if (error || !isEmpty(errors)) {
            hasErrors = true;
            this.$popups.error({ text: error || Object.values(errors).join('\n') });
          }
        });

        if (!hasErrors) {
          this.$popups.success({
            text: this.$t(`modals.createPbehavior.success.${remove ? 'remove' : 'create'}`),
          });
        }
      }
    },

    async createDowntimePbehavior(entities, payload) {
      const response = await this.bulkCreateEntityPbehaviors({
        data: entities.map(
          entity =>
            createDowntimePbehavior({
              entity,
              ...pick(payload, ['comment', 'reason', 'type', 'prefix', 'origin']),
            }),
          [],
        ),
      });

      this.showPbehaviorResponsePopups(response);
    },

    async removeDowntimePbehavior(entities, origin = PBEHAVIOR_ORIGINS.serviceWeather) {
      const response = await this.bulkRemoveEntityPbehaviors({
        data: entities.map(({ _id: id }) => ({
          origin,
          entity: id,
        })),
      });

      this.showPbehaviorResponsePopups(response, true);
    },
  },
};
