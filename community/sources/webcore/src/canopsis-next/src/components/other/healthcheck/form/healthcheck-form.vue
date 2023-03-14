<template lang="pug">
  v-layout(column)
    c-information-block(:title="$t('healthcheck.queueLimit')")
      template(#subtitle="") {{ $t('healthcheck.notifyUsersQueueLimit') }}
      c-enabled-limit-field(
        v-field="form.queue",
        :label="$t('healthcheck.defineQueueLimit')"
      )
    c-information-block(:title="$t('healthcheck.numberOfInstances')")
      template(#subtitle="") {{ $t('healthcheck.notifyUsersNumberOfInstances') }}
      healthcheck-engine-instance-field(
        v-for="engineName in engineNames",
        v-field="form[engineName]",
        :name="engineName",
        :key="engineName",
        :label="$t(`healthcheck.nodes.${engineName}.name`)"
      )
</template>

<script>
import { HEALTHCHECK_ENGINES_NAMES } from '@/constants';

import HealthcheckEngineInstanceField from './partials/healthcheck-engine-instance-field.vue';

export default {
  inject: ['$validator'],
  components: { HealthcheckEngineInstanceField },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
  },
  computed: {
    engineNames() {
      return Object.values(HEALTHCHECK_ENGINES_NAMES);
    },
  },
};
</script>
