import { Link } from "@tanstack/react-router";

import { ThemeToggle } from "./theme-togle";
import {
	Sidebar,
	SidebarContent,
	SidebarFooter,
	SidebarGroup,
	SidebarGroupContent,
	SidebarGroupLabel,
	SidebarHeader,
	SidebarMenu,
	SidebarMenuItem,
} from "./ui";

export const AppSidebar = () => {
	return (
		<Sidebar className="h-full flex flex-col items-start">
			<SidebarHeader>
				<Link to="/">
					<h1 className="text-2xl font-semibold">My App</h1>
				</Link>
			</SidebarHeader>
			<SidebarContent>
				<SidebarGroup>
					<SidebarGroupLabel>Application</SidebarGroupLabel>
					<SidebarGroupContent>
						<SidebarMenu title="Main Menu" />
						<Link to="/dashboard">
							<SidebarMenuItem className="p-2">Todos</SidebarMenuItem>
						</Link>
						<SidebarMenu />
					</SidebarGroupContent>
				</SidebarGroup>
			</SidebarContent>

			<SidebarFooter className="mt-auto w-full flex flex-row justify-between items-center">
				<h1 className="text-2xl">Hi, Akshay</h1>
				<ThemeToggle />
			</SidebarFooter>
		</Sidebar>
	);
};
