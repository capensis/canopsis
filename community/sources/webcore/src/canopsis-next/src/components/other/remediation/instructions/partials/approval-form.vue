<template>
  <v-form @submit.prevent="submitDismiss">
    <v-row>
      <v-col v-if="dismissing" cols="12">
        <c-description-field
          v-model="comment"
          :label="$t('common.comment')"
          name="dismiss_comment"
          required
        />
      </v-col>

      <v-col class="d-flex justify-end" cols="12">
        <template v-if="!dismissing">
          <v-btn
            :disabled="disabled || submitting"
            class="mr-2"
            color="warning"
            @click="showDismissComment"
          >
            {{ $t('common.dismiss') }}
          </v-btn>
          <v-btn
            :disabled="disabled || submitting"
            :loading="submitting"
            color="primary"
            @click="$emit('approve')"
          >
            {{ $t('common.approve') }}
          </v-btn>
        </template>
        <template v-else>
          <v-btn
            :disabled="submitting"
            class="mr-2"
            variant="text"
            @click="cancelDismiss"
          >
            {{ $t('common.cancel') }}
          </v-btn>
          <v-btn
            :disabled="disabled || submitting"
            :loading="submitting"
            color="warning"
            type="submit"
          >
            {{ $t('common.dismiss') }}
            <v-icon color="error" end>
              cancel
            </v-icon>
          </v-btn>
        </template>
      </v-col>
    </v-row>
  </v-form>
</template>

<script setup>
import { ref, watch } from 'vue';

import { useForm } from '@/hooks/form';

defineProps({
  submitting: {
    type: Boolean,
    default: false,
  },
  disabled: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(['approve', 'dismiss']);

const dismissing = ref(false);
const form = useForm({
  initialValues: {
    comment: '',
  },
});
const comment = form.model('comment');

const closeDismissComment = () => dismissing.value = false;

const showDismissComment = () => dismissing.value = true;

const cancelDismiss = () => {
  comment.value = '';
  closeDismissComment();
};

const submitDismiss = form.handleSubmit(() => {
  emit('dismiss', comment.value);
});

watch(dismissing, () => {
  form.resetErrors();
});
</script>
