import { isEmpty, pick } from 'lodash';

import { PBEHAVIOR_ORIGINS } from '@/constants';

import { createDowntimePbehavior } from '@/helpers/entities/pbehavior/form';

import { entitiesPbehaviorMixin } from '@/mixins/entities/pbehavior';

export const widgetActionsPanelCommonMixin = {
  mixins: [entitiesPbehaviorMixin],
  methods: {
    showPbehaviorResponseErrorPopups(response) {
      if (response?.length) {
        response.forEach(({ error, errors }) => {
          if (error || !isEmpty(errors)) {
            this.$popups.error({ text: error || Object.values(errors).join('\n') });
          }
        });
      }
    },

    async createDowntimePbehavior(entities, payload) {
      const response = await this.createEntityPbehaviors({
        data: entities.map(entity => createDowntimePbehavior({
          entity,
          ...pick(payload, ['comment', 'reason', 'type', 'prefix', 'origin']),
        }), []),
      });

      this.showPbehaviorResponseErrorPopups(response);
    },

    async removeDowntimePbehavior(entities) {
      const response = await this.removeEntityPbehaviors({
        data: entities.map(({ _id: id }) => ({
          entity: id,
          origin: PBEHAVIOR_ORIGINS.serviceWeather,
        })),
      });

      this.showPbehaviorResponseErrorPopups(response);
    },
  },
};
