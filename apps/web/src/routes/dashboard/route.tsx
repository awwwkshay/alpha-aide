import { createFileRoute, Outlet } from "@tanstack/react-router";

import { AppSidebar, SidebarProvider, SidebarTrigger } from "@/components";

export const Route = createFileRoute("/dashboard")({
	component: DashboardLayout,
});

function DashboardLayout() {
	return (
		<div className="dashboard-layout">
			<SidebarProvider>
				<AppSidebar />
				<section className="w-full h-screen">
					<SidebarTrigger />
					<main>
						<Outlet />
					</main>
				</section>
			</SidebarProvider>
		</div>
	);
}
