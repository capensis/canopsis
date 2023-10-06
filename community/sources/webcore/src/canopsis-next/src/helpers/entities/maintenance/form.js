import { DEFAULT_BROADCAST_MESSAGE_COLOR } from '@/constants';

/**
 * @typedef {Object} Maintenance
 * @property {string} color
 * @property {boolean} [enabled]
 * @property {string} message
 */

/**
 * @typedef {Maintenance} MaintenanceForm
 */

/**
 * Convert maintenance to form
 *
 * @param {Maintenance} maintenance
 * @returns {MaintenanceForm}
 */
export const maintenanceToForm = (maintenance = {}) => ({
  color: maintenance.color ?? DEFAULT_BROADCAST_MESSAGE_COLOR,
  message: maintenance.message ?? '',
});
