import { VUETIFY_ANIMATION_DELAY } from '@/config';

/**
 * Wait a vuetify animation
 *
 * @return {Promise}
 */
export const waitVuetifyAnimation = () => new Promise(resolve => setTimeout(resolve, VUETIFY_ANIMATION_DELAY));
