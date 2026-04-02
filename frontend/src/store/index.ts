import { configureStore } from '@reduxjs/toolkit';
import authReducer from './authSlice';
import workflowReducer from './workflowSlice';
import notificationReducer from './notificationSlice';

export const store = configureStore({
  reducer: {
    auth: authReducer,
    workflow: workflowReducer,
    notification: notificationReducer,
  },
});

export type RootState = ReturnType<typeof store.getState>
export type AppDispatch = typeof store.dispatch
