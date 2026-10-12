<template>
  <c-layout column>
    <c-form-block :title="$t('healthcheck.queueLimit')">
      <template #subtitle>
        {{ $t('healthcheck.notifyUsersQueueLimit') }}
      </template>
      <c-form-block-row :label="$t('healthcheck.defineQueueLimit')">
        <healthcheck-enabled-limit-field
          :value="form.queue"
          :label="$t('healthcheck.defineQueueLimit')"
          :limit-errors="getErrors('queue.limit')"
          :min="1"
          name="queue"
          @input="updateField('queue', $event)"
        />
      </c-form-block-row>
    </c-form-block>

    <c-form-block :title="$t('healthcheck.messagesLimit')" class="mt-4">
      <template #subtitle>
        {{ $t('healthcheck.notifyUsersMessagesLimit') }}
      </template>
      <c-form-block-row :label="$t('healthcheck.defineMessageLimit')">
        <healthcheck-enabled-limit-field
          :value="form.messages"
          :label="$t('healthcheck.defineMessageLimit')"
          :limit-errors="getErrors('messages.limit')"
          :min="1"
          name="messages"
          @input="updateField('messages', $event)"
        />
      </c-form-block-row>
    </c-form-block>

    <c-form-block :title="$t('healthcheck.numberOfInstances')" class="mt-4">
      <template #subtitle>
        {{ $t('healthcheck.notifyUsersNumberOfInstances') }}
      </template>
      <c-form-block-row
        v-for="engineName in engineNames"
        :key="engineName"
        :label="$t(`healthcheck.nodes.${engineName}.name`)"
      >
        <healthcheck-engine-instance-field
          :value="form[engineName]"
          :name="engineName"
          :label="$t(`healthcheck.nodes.${engineName}.name`)"
          :minimal-errors="getErrors(`${engineName}.minimal`)"
          :optimal-errors="getErrors(`${engineName}.optimal`)"
          @input="updateField(engineName, $event)"
        />
      </c-form-block-row>
    </c-form-block>
  </c-layout>
</template>

<script>
import { HEALTHCHECK_ENGINES_NAMES } from '@/constants';

import HealthcheckEnabledLimitField from './fields/healthcheck-enabled-limit-field.vue';
import HealthcheckEngineInstanceField from './fields/healthcheck-engine-instance-field.vue';

export default {
  components: {
    HealthcheckEnabledLimitField,
    HealthcheckEngineInstanceField,
  },

  model: {
    prop: 'form',
    event: 'input',
  },

  props: {
    form: {
      type: Object,
      required: true,
    },
    errors: {
      type: Object,
      default: () => ({}),
    },
  },

  emits: ['input'],

  computed: {
    engineNames() {
      return Object.values(HEALTHCHECK_ENGINES_NAMES);
    },
  },

  methods: {
    updateField(field, value) {
      this.$emit('input', {
        ...this.form,
        [field]: value,
      });
    },

    getErrors(field) {
      return this.errors[field] || [];
    },
  },
};
</script>
