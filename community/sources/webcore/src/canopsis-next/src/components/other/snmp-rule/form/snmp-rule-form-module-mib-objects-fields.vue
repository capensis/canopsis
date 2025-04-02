<template>
  <v-layout wrap>
    <v-flex xs12>
      <v-menu
        :items="items"
        max-height="200"
        offset-y
      >
        <template #activator="{ on }">
          <v-text-field
            v-field="form.value"
            v-validate="valueRules"
            :label="$t('snmpRule.moduleMibObjects')"
            :name="valueName"
            :error-messages="errors.collect(valueName)"
            class="vars-input pt-0"
            hide-details
            v-on="on"
          >
            <template
              v-if="large"
              #append=""
            >
              <v-btn
                :class="{ active: isVisible }"
                icon
                @click.stop="toggleVisibility"
              >
                <v-icon>attach_file</v-icon>
              </v-btn>
            </template>
          </v-text-field>
        </template>
        <v-list>
          <v-list-item
            v-for="(item, index) in items"
            :key="index"
            @click="updateSelectableInput(item)"
          >
            <v-list-item-title>{{ item }}</v-list-item-title>
          </v-list-item>
        </v-list>
      </v-menu>
    </v-flex>
    <v-expand-transition v-if="large">
      <v-flex v-show="isVisible" xs12>
        <v-text-field
          v-field="form.regex"
          :label="$t('snmpRule.regex')"
          hide-details
        />
        <v-text-field
          v-field="form.formatter"
          :label="$t('snmpRule.formatter')"
          hide-details
        />
      </v-flex>
    </v-expand-transition>
  </v-layout>
</template>

<script>
import { computed, ref } from 'vue';

import { useModelField } from '@/hooks/form/model-field';

export default {
  inject: ['$validator'],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
    items: {
      type: Array,
      default: () => [],
    },
    large: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      default: '',
    },
    required: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const { updateField } = useModelField(props, emit);

    const isVisible = ref(!!(props.form.regex || props.form.format));

    const valueRules = computed(() => ({ required: props.required }));
    const valueName = computed(() => (props.name ? `${props.name}.value` : 'value'));

    const toggleVisibility = () => isVisible.value = !isVisible.value;
    const updateSelectableInput = item => updateField('value', `${props.form.value || ''}{{ ${item} }}`);

    return {
      isVisible,

      valueRules,
      valueName,

      toggleVisibility,
      updateSelectableInput,
    };
  },
};
</script>

<style lang="scss" scoped>
.v-btn.active {
  &:hover:before {
    opacity: .16;
  }

  &:before {
    background-color: currentColor;
  }
}

.vars-input ::v-deep .v-input__slot {
  height: 56px;
}
</style>
