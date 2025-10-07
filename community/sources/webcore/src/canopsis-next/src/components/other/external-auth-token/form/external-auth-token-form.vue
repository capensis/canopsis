<template>
  <v-layout column>
    <c-name-field
      v-field="form.name"
      :label="$t('externalAuthToken.tokenName')"
      name="name"
      autofocus
      required
    />
    <c-description-field v-field="form.description" />
    <request-form
      v-field="form.request"
      name="request"
    />
    <c-information-block :title="$t('common.token')">
      <v-layout class="gap-2">
        <v-flex xs6>
          <c-enabled-field
            v-field="form.allow_variables"
            :label="$t('externalAuthToken.allowVariables')"
          />
        </v-flex>
        <v-flex xs6>
          <v-text-field
            v-field="form.template"
            v-validate="'required'"
            :label="$t('common.token')"
            :error-messages="errors.collect('token')"
            name="token"
          >
            <template #append="">
              <c-help-icon
                :text="$t('externalAuthToken.tokenExpirationHelpText')"
                icon="help"
                color="grey darken-1"
                top
              />
            </template>
          </v-text-field>
        </v-flex>
      </v-layout>
    </c-information-block>
    <c-information-block :title="$t('externalAuthToken.tokenExpirationTime')">
      <c-duration-field
        v-field="form.expiration_duration"
        required
        same-width
      />
    </c-information-block>
  </v-layout>
</template>

<script>
import RequestForm from '@/components/forms/request/request-form.vue';

export default {
  inject: ['$validator'],
  components: { RequestForm },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    isNew: {
      type: Boolean,
      default: false,
    },
    fromConfig: {
      type: Boolean,
      default: false,
    },
  },
  setup() {
    const nameRegex = /^[a-z_][\w_]*$/i;

    return {
      nameRegex,
    };
  },
};
</script>
