<template>
  <v-row density="comfortable">
    <v-col cols="12">
      <v-menu
        v-model="menu"
        :items="items"
        max-height="200"
        location="bottom"
      >
        <template #activator="{ props: menuProps }">
          <c-text-field
            v-bind="menuProps"
            v-model="value"
            :label="label || t('snmpRule.moduleMibObjects')"
            :name="valueName"
            :required="required"
            class="vars-input pt-0"
            hide-details
          >
            <template v-if="large" #append-inner>
              <v-btn
                :class="{ active: isVisible }"
                size="small"
                variant="text"
                icon
                @click.stop="toggleVisibility"
              >
                <v-icon>attach_file</v-icon>
              </v-btn>
            </template>
          </c-text-field>
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
    </v-col>

    <v-expand-transition v-if="large">
      <v-col v-show="isVisible" cols="12">
        <c-text-field
          v-model="regex"
          :label="t('snmpRule.regex')"
          :name="regexName"
          hide-details
        />
        <c-text-field
          v-model="formatter"
          :label="t('snmpRule.formatter')"
          :name="formatterName"
        />
      </v-col>
    </v-expand-transition>
  </v-row>
</template>

<script>
import { computed, ref } from 'vue';

import { useI18n } from '@/hooks/i18n';
import { useModelField } from '@/hooks/form/model-field';

export default {
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
    label: {
      type: String,
      default: '',
    },
    required: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:form'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const { updateField } = useModelField(props, emit);

    const menu = ref(false);
    const isVisible = ref(!!(props.form.regex || props.form.formatter));

    const valueName = computed(() => (props.name ? `${props.name}.value` : 'value'));
    const regexName = computed(() => (props.name ? `${props.name}.regex` : 'regex'));
    const formatterName = computed(() => (props.name ? `${props.name}.formatter` : 'formatter'));

    const value = computed({
      get: () => props.form.value,
      set: nextValue => updateField('value', nextValue),
    });

    const regex = computed({
      get: () => props.form.regex,
      set: nextValue => updateField('regex', nextValue),
    });

    const formatter = computed({
      get: () => props.form.formatter,
      set: nextValue => updateField('formatter', nextValue),
    });

    const toggleVisibility = () => {
      isVisible.value = !isVisible.value;
    };

    const updateSelectableInput = (item) => {
      updateField('value', `${props.form.value || ''}{{ ${item} }}`);
      menu.value = false;
    };

    return {
      t,
      menu,
      isVisible,
      valueName,
      regexName,
      formatterName,
      value,
      regex,
      formatter,
      toggleVisibility,
      updateSelectableInput,
    };
  },
};
</script>

<style lang="scss" scoped>
.v-btn.active {
  &:hover::before {
    opacity: 0.16;
  }

  &::before {
    background-color: currentColor;
  }
}

.vars-input :deep(.v-input__slot) {
  height: 56px;
}

.vars-input :deep(.v-field__append-inner) {
  padding-top: 4px;
}
</style>
