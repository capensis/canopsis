import { createNamespacedHelpers } from 'vuex';

import { MAX_LIMIT } from '@/constants';

const { mapActions: mapServiceActions } = createNamespacedHelpers('service');

export const entityVariablesMixin = {
  data() {
    return {
      infos: [],
    };
  },
  mounted() {
    this.fetchInfos();
  },
  computed: {
    infosSubVariables() {
      return [
        {
          text: this.$t('common.value'),
          value: 'value',
        },
        {
          text: this.$t('common.description'),
          value: 'description',
        },
      ];
    },

    infosVariables() {
      return this.infos.map(({ value }) => ({
        text: value,
        value,
        variables: this.infosSubVariables,
      }));
    },

    stateVariables() {
      return [
        {
          text: this.$t('common.timestamp'),
          value: 't',
        },
        {
          text: this.$t('common.value'),
          value: 'val',
        },
      ];
    },

    statusVariables() {
      return [
        {
          text: this.$t('common.value'),
          value: 'val',
        },
      ];
    },

    variables() {
      return [
        {
          text: this.$t('common.id'),
          value: 'entity._id',
        },
        {
          text: this.$t('common.name'),
          value: 'entity.name',
        },
        {
          text: this.$t('common.infos'),
          value: 'entity.infos',
          variables: this.infosVariables,
        },
        {
          text: this.$t('common.connector'),
          value: 'entity.connector',
        },
        {
          text: this.$t('common.connectorName'),
          value: 'entity.connector_name',
        },
        {
          text: this.$t('common.component'),
          value: 'entity.component',
        },
        {
          text: this.$t('common.resource'),
          value: 'entity.resource',
        },
        {
          text: this.$t('common.state'),
          value: 'entity.state',
          variables: this.stateVariables,
        },
        {
          text: this.$t('common.status'),
          value: 'entity.status',
          variables: this.statusVariables,
        },
        {
          text: this.$t('common.snooze'),
          value: 'entity.snooze',
        },
        {
          text: this.$t('common.ack'),
          value: 'entity.ack',
        },
        {
          text: this.$t('common.updated'),
          value: 'entity.last_update_date',
        },
        {
          text: this.$t('common.impactLevel'),
          value: 'entity.impact_level',
        },
        {
          text: this.$t('common.impactState'),
          value: 'entity.impact_state',
        },
        {
          text: this.$t('common.category'),
          value: 'entity.category.name',
        },
        {
          text: this.$tc('common.link', 2),
          value: 'entity.links',
        },
      ];
    },
  },
  methods: {
    ...mapServiceActions({ fetchEntityInfosKeysWithoutStore: 'fetchInfosKeysWithoutStore' }),

    async fetchInfos() {
      const { data: infos } = await this.fetchEntityInfosKeysWithoutStore({
        params: { limit: MAX_LIMIT },
      });

      this.infos = infos;
    },
  },
};
