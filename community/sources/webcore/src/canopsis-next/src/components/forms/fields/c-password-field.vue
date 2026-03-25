<template>
  <v-layout v-if="replaceable" class="gap-4" column>
    <v-label :aria-required="required" class="c-password-field__label">
      {{ label }}
    </v-label>
    <v-layout v-if="shownField" class="gap-3" align-center>
      <v-text-field
        v-field="value"
        v-validate="rules"
        v-bind="$attrs"
        :placeholder="placeholder || $t('common.password')"
        :error-messages="errors.collect(name)"
        :name="name"
        :type="shownPassword ? 'text' : 'password'"
        :append-icon="appendIcon"
        :clearable="false"
        @click:append="toggleShownPassword"
      />
      <c-action-btn
        type="delete"
        top
        small
        @click="toggleShownField"
      />
    </v-layout>
    <div v-if="!shownField">
      <v-btn color="primary" outlined @click="toggleShownField">
        {{ buttonLabel || `${$t('common.replace')} ${label}` }}
      </v-btn>
    </div>
  </v-layout>
  <v-text-field
    v-else
    v-field="value"
    v-validate="rules"
    v-bind="$attrs"
    :label="label || $t('common.password')"
    :error-messages="errors.collect(name)"
    :name="name"
    :type="shownPassword ? 'text' : 'password'"
    :append-icon="appendIcon"
    :clearable="false"
    @click:append="toggleShownPassword"
  />
</template>

<script>
import { computed, ref } from 'vue';

export default {
  inject: ['$validator'],
  inheritAttrs: false,
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      default: '',
    },
    label: {
      type: String,
      default: '',
    },
    placeholder: {
      type: String,
      default: '',
    },
    buttonLabel: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'password',
    },
    required: {
      type: Boolean,
      default: false,
    },
    replaceable: {
      type: Boolean,
      default: false,
    },
    visibility: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const shownPassword = ref(props.replaceable);
    const shownField = ref(false);

    const rules = computed(() => ({
      required: props.required,
    }));

    const appendIcon = computed(() => {
      if (!props.visibility) {
        return null;
      }

      return shownPassword.value ? 'visibility' : 'visibility_off';
    });

    const toggleShownPassword = () => shownPassword.value = !shownPassword.value;
    const toggleShownField = () => shownField.value = !shownField.value;

    return {
      shownPassword,
      shownField,
      rules,
      appendIcon,
      toggleShownPassword,
      toggleShownField,
    };
  },
};
</script>

<style lang="scss" scoped>
.c-password-field__label {
  transform: scale(0.75);
  transform-origin: left;
}
</style>
