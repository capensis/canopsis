<template>
  <v-card>
    <v-card-text>
      <c-name-field
        v-field="form.name"
        :name="nameFieldName"
        required
      >
        <template #append-outer>
          <c-action-btn
            type="delete"
            @click="$emit('remove', form)"
          />
        </template>
      </c-name-field>
      <v-layout class="gap-4">
        <c-enabled-field
          v-field="form.copy_from_children"
          :label="$t('metaAlarmRule.copyFromLastChild')"
          class="mr-4 pt-4"
        >
          <template #append>
            <c-help-icon
              :text="$t('metaAlarmRule.copyFromLastChildTooltip')"
              icon="help"
              top
            />
          </template>
        </c-enabled-field>
        <v-flex>
          <c-mixed-field
            v-field="form.value"
            :label="$t('common.value')"
          />
        </v-flex>
      </v-layout>
      <v-text-field
        v-field="form.description"
        :label="$t('common.description')"
      />
    </v-card-text>
  </v-card>
</template>

<script>
import { computed } from 'vue';

export default {
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    name: {
      type: String,
      default: '',
    },
  },
  setup(props) {
    const getFieldName = field => `${props.name}.${field}`;

    const nameFieldName = computed(() => getFieldName('name'));

    return {
      nameFieldName,
    };
  },
};
</script>
