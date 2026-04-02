import { createSlice, PayloadAction } from '@reduxjs/toolkit';

interface WorkflowState {
  definitions: any[]
  currentDefinition: any | null
}

const initialState: WorkflowState = {
  definitions: [],
  currentDefinition: null,
};

const workflowSlice = createSlice({
  name: 'workflow',
  initialState,
  reducers: {
    setDefinitions: (state, action: PayloadAction<any[]>) => {
      state.definitions = action.payload;
    },
    setCurrentDefinition: (state, action: PayloadAction<any>) => {
      state.currentDefinition = action.payload;
    },
  },
});

export const { setDefinitions, setCurrentDefinition } = workflowSlice.actions;
export default workflowSlice.reducer;
