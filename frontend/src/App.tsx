import { Routes, Route } from 'react-router-dom';
import { Layout } from './components/common/Layout/Layout';
import { Login } from './pages/Login';
import { Dashboard } from './pages/Dashboard';
import { CompanyList } from './pages/org/CompanyList';
import { DepartmentList } from './pages/org/DepartmentList';
import { EmployeeList } from './pages/org/EmployeeList';
import { PositionList } from './pages/admin/PositionList';
import { SystemSettingsPage } from './pages/admin/SystemSettings';
import { WorkflowList } from './pages/workflow/DefinitionList';
import { WorkflowDesigner } from './pages/workflow/Designer/WorkflowDesigner';
import { WorkflowApply } from './pages/workflow/WorkflowApply';
import { MyApplications } from './pages/workflow/MyApplications';
import { MyTasks } from './pages/approval/MyTasks';
import { HandledTasks } from './pages/approval/HandledTasks';
import { NotificationList } from './pages/notification';

function App() {
  return (
    <Routes>
      <Route path="/login" element={<Login />} />
      <Route path="/" element={<Layout />}>
        <Route index element={<Dashboard />} />
        <Route path="companies" element={<CompanyList />} />
        <Route path="departments" element={<DepartmentList />} />
        <Route path="employees" element={<EmployeeList />} />
        <Route path="positions" element={<PositionList />} />
        <Route path="system-settings" element={<SystemSettingsPage />} />
        <Route path="workflows" element={<WorkflowList />} />
        <Route path="workflows/designer/:id?" element={<WorkflowDesigner />} />
        <Route path="workflows/apply/:id" element={<WorkflowApply />} />
        <Route path="my-applications" element={<MyApplications />} />
        <Route path="tasks/pending" element={<MyTasks />} />
        <Route path="tasks/handled" element={<HandledTasks />} />
        <Route path="notifications" element={<NotificationList />} />
      </Route>
    </Routes>
  );
}

export default App;
