<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper :close="cancel">
      <template #title>
        <span>{{ $t('modals.pbehaviorRecurrentChangesConfirmation.title') }}</span>
      </template>

      <template #text>
        <c-radio-group v-model="type" hide-details>
          <v-radio
            :value="PBEHAVIOR_PLANNING_EVENT_CHANGING_TYPES.selected"
            :label="$t('modals.pbehaviorRecurrentChangesConfirmation.fields.selected')"
            color="primary"
          />
          <v-radio
            :value="PBEHAVIOR_PLANNING_EVENT_CHANGING_TYPES.all"
            :label="$t('modals.pbehaviorRecurrentChangesConfirmation.fields.all')"
            color="primary"
          />
        </c-radio-group>
      </template>

      <template #actions>
        <v-btn variant="text" @click="cancel">
          {{ $t('common.cancel') }}
        </v-btn>

        <v-btn color="primary" type="submit">
          {{ $t('common.submit') }}
        </v-btn>
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script>
import { PBEHAVIOR_PLANNING_EVENT_CHANGING_TYPES } from '@/constants';

import { useModals } from '@/hooks/modals';

export default {
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  data() {
    return {
      type: PBEHAVIOR_PLANNING_EVENT_CHANGING_TYPES.selected,
    };
  },

  setup() {
    const modals = useModals();

    return {
      modals,
      PBEHAVIOR_PLANNING_EVENT_CHANGING_TYPES,
    };
  },

  computed: {
    config() {
      return this.modal.config ?? {};
    },
  },

  methods: {
    submit() {
      this.config.action?.(this.type);

      this.modals.hide({ id: this.modal.id });
    },

    cancel() {
      this.config.cancel?.();

      this.modals.hide({ id: this.modal.id });
    },
  },
};
</script>
