import { WEBHOOK_EXECUTION_STATUSES } from '@/constants';

/**
 * Check webhook execution status is waiting
 *
 * @param {Object} execution
 * @returns {boolean}
 */
export const isWebhookExecutionWaiting = execution => execution?.status === WEBHOOK_EXECUTION_STATUSES.waiting;

/**
 * Check webhook execution status is running
 *
 * @param {Object} execution
 * @returns {boolean}
 */
export const isWebhookExecutionRunning = execution => execution?.status === WEBHOOK_EXECUTION_STATUSES.running;

/**
 * Check webhook execution status is succeeded
 *
 * @param {Object} execution
 * @returns {boolean}
 */
export const isWebhookExecutionSucceeded = execution => execution?.status === WEBHOOK_EXECUTION_STATUSES.succeeded;

/**
 * Check webhook execution status is failed
 *
 * @param {Object} execution
 * @returns {boolean}
 */
export const isWebhookExecutionFailed = execution => execution?.status === WEBHOOK_EXECUTION_STATUSES.failed;

/**
 * Check webhook execution is finished
 *
 * @param {Object} execution
 * @returns {boolean}
 */
export const isWebhookExecutionFinished = execution => !!execution
  && (isWebhookExecutionSucceeded(execution) || isWebhookExecutionFailed(execution));
