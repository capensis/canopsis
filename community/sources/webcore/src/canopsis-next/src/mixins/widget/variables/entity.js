import { createNamespacedHelpers } from 'vuex';

import { MAX_LIMIT, ENTITY_TEMPLATE_FIELDS, ENTITY_FIELDS_TO_LABELS_KEYS } from '@/constants';

import { variablesMixin } from './common';

const { mapActions: mapServiceActions } = createNamespacedHelpers('service');

export const entityVariablesMixin = {
  mixins: [variablesMixin],
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

    snoozeVariables() {
      return [
        this.alarmStepTimestampVariable,
        this.alarmStepAuthorVariable,
        this.alarmStepMessageVariable,
      ];
    },

    entityVariables() {
      return [
        { value: ENTITY_TEMPLATE_FIELDS.id },
        { value: ENTITY_TEMPLATE_FIELDS.name },
        {
          value: ENTITY_TEMPLATE_FIELDS.infos,
          variables: this.infosVariables,
        },
        { value: ENTITY_TEMPLATE_FIELDS.connector },
        { value: ENTITY_TEMPLATE_FIELDS.connectorName },
        { value: ENTITY_TEMPLATE_FIELDS.component },
        { value: ENTITY_TEMPLATE_FIELDS.resource },
        {
          value: ENTITY_TEMPLATE_FIELDS.state,
          variables: this.stateVariables,
        },
        {
          value: ENTITY_TEMPLATE_FIELDS.status,
          variables: this.statusVariables,
        },
        {
          value: ENTITY_TEMPLATE_FIELDS.snooze,
          variables: this.snoozeVariables,
        },
        {
          value: ENTITY_TEMPLATE_FIELDS.ack,
          variables: this.ackVariables,
        },
        { value: ENTITY_TEMPLATE_FIELDS.lastUpdateDate },
        { value: ENTITY_TEMPLATE_FIELDS.impactLevel },
        { value: ENTITY_TEMPLATE_FIELDS.impactState },
        { value: ENTITY_TEMPLATE_FIELDS.categoryName },
        {
          value: ENTITY_TEMPLATE_FIELDS.pbehaviorInfo,
          variables: this.pbehaviorInfoVariables,
        },
        { value: ENTITY_TEMPLATE_FIELDS.alarmCreationDate },
        {
          value: ENTITY_TEMPLATE_FIELDS.ticket,
          variables: this.ticketVariables,
        },
        { value: ENTITY_TEMPLATE_FIELDS.statsOk },
        { value: ENTITY_TEMPLATE_FIELDS.statsKo },
        { value: ENTITY_TEMPLATE_FIELDS.alarmDisplayName },
        { value: ENTITY_TEMPLATE_FIELDS.links },
      ].map(variable => ({
        ...variable,

        text: this.$tc(ENTITY_FIELDS_TO_LABELS_KEYS[variable.value], 2),
      }));
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
