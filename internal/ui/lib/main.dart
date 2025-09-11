// not a mobile dev, flutter is newer to me, so learning with chatgpt as we go:)
// as this scales a little more, will need to break apart and modularize as much as possible, here future me: https://docs.flutter.dev/app-architecture/guide
import 'package:flutter/material.dart';

void main() {
  runApp(const MainApp());
}

class MainApp extends StatefulWidget {
  const MainApp({super.key});

  @override
  State<MainApp> createState() => _MainAppState();
}

class _MainAppState extends State<MainApp> {
  int _selectedIndex = 0;

  static const accentColor = Colors.tealAccent;

  late final List<Widget> _pages = [
    const SearchPage(),
    const StarredPage(),
    const ProfilePage(),
  ];

  void _onItemTapped(int index) {
    setState(() {
      _selectedIndex = index;
    });
  }

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      debugShowCheckedModeBanner: false,
      theme: ThemeData(
        brightness: Brightness.dark,
        useMaterial3: true,
        scaffoldBackgroundColor: const Color(0xFF121212),
        appBarTheme: const AppBarTheme(
	  backgroundColor: Colors.transparent,
          elevation: 0,
        ),
        textTheme: const TextTheme(
          bodyMedium: TextStyle(color: Colors.white70, fontSize: 18),
        ),
        colorScheme: ColorScheme.fromSeed(
          seedColor: accentColor,
          brightness: Brightness.dark,
        ),
      ),
      home: Scaffold(
        body: AnimatedSwitcher(
          duration: const Duration(milliseconds: 300),
          child: _pages[_selectedIndex],
        ),
        bottomNavigationBar: NavigationBar(
          backgroundColor: const Color(0xFF1F1F1F),
          selectedIndex: _selectedIndex,
          onDestinationSelected: _onItemTapped,
          destinations: const [
            NavigationDestination(icon: Icon(Icons.search), label: "Search"),
            NavigationDestination(icon: Icon(Icons.star), label: "Saved"),
            NavigationDestination(icon: Icon(Icons.person), label: "Profile"),
          ],
        ),
      ),
    );
  }
}

// SEARCH PAGE / HOME PAGE
class SearchPage extends StatefulWidget {
  const SearchPage({super.key});

  @override
  State<SearchPage> createState() => _SearchPageState();
}

class _SearchPageState extends State<SearchPage> {
  final TextEditingController _titleController = TextEditingController();
  final TextEditingController _zipController = TextEditingController();
  final TextEditingController _radiusController = TextEditingController();

  @override
  void dispose() {
    _titleController.dispose();
    _zipController.dispose();
    _radiusController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return SafeArea(
      child: CustomScrollView(
        slivers: [
          // nav bar with search fields
          SliverAppBar(
            pinned: true,
            floating: true,
            expandedHeight: 180, // slightly taller to prevent overflow
            flexibleSpace: FlexibleSpaceBar(
              background: Padding(
                padding: const EdgeInsets.fromLTRB(16, 40, 16, 8),
                child: SingleChildScrollView(
                  physics: const NeverScrollableScrollPhysics(), // prevents nested scroll conflicts
                  child: Column(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      // title input
                      TextField(
                        controller: _titleController,
                        decoration: InputDecoration(
                          hintText: "Job Title",
                          prefixIcon: const Icon(Icons.work, color: Colors.white70),
                          filled: true,
                          fillColor: Colors.grey[850],
                          border: OutlineInputBorder(
                            borderRadius: BorderRadius.circular(12),
                            borderSide: BorderSide.none,
                          ),
                        ),
                      ),
                      const SizedBox(height: 8),
                      // zip + radius input
                      Row(
                        children: [
                          Expanded(
                            child: TextField(
                              controller: _zipController,
                              keyboardType: TextInputType.number,
                              decoration: InputDecoration(
                                hintText: "Zip Code",
                                prefixIcon:
                                    const Icon(Icons.location_on, color: Colors.white70),
                                filled: true,
                                fillColor: Colors.grey[850],
                                border: OutlineInputBorder(
                                  borderRadius: BorderRadius.circular(12),
                                  borderSide: BorderSide.none,
                                ),
                              ),
                            ),
                          ),
                          const SizedBox(width: 12),
                          Expanded(
                            child: TextField(
                              controller: _radiusController,
                              keyboardType: TextInputType.number,
                              decoration: InputDecoration(
                                hintText: "Radius (miles)",
                                prefixIcon: const Icon(Icons.swap_vert,
                                    color: Colors.white70),
                                filled: true,
                                fillColor: Colors.grey[850],
                                border: OutlineInputBorder(
                                  borderRadius: BorderRadius.circular(12),
                                  borderSide: BorderSide.none,
                                ),
                              ),
                            ),
                          ),
                        ],
                      ),
                    ],
                  ),
                ),
              ),
            ),
          ),

        // listings
          SliverList(
            delegate: SliverChildBuilderDelegate(
              (context, index) {
                return JobCard(
                  title: "Job Title $index",
                  company: "Company $index",
                  location: "Location $index",
                );
              },
              childCount: 25, // filler
            ),
          ),
        ],
      ),
    );
  }
}

// STARRED PAGE
class StarredPage extends StatelessWidget {
  const StarredPage({super.key});

  @override
  Widget build(BuildContext context) {
    return SafeArea(
      child: ListView.builder(
        padding: const EdgeInsets.all(16),
        itemCount: 10, // example count
        itemBuilder: (context, index) {
          return JobCard(
            title: "Saved Job $index",
            company: "Company $index",
            location: "Location $index",
            isSaved: true,
          );
        },
      ),
    );
  }
}

// PROFILE PAGE
class ProfilePage extends StatelessWidget {
  const ProfilePage({super.key});

  @override
  Widget build(BuildContext context) {
    return SafeArea(
      child: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: const [
            CircleAvatar(
              radius: 48,
              backgroundColor: Colors.tealAccent,
              child: Icon(Icons.person, size: 48, color: Colors.black87),
            ),
            SizedBox(height: 16),
            Text(
              "User Name",
              style: TextStyle(fontSize: 22, fontWeight: FontWeight.bold),
            ),
            SizedBox(height: 8),
            Text("user@email.com", style: TextStyle(color: Colors.white70)),
          ],
        ),
      ),
    );
  }
}

// JOB CARD WIDGET
class JobCard extends StatelessWidget {
  final String title;
  final String company;
  final String location;
  final bool isSaved;

  const JobCard({
    super.key,
    required this.title,
    required this.company,
    required this.location,
    this.isSaved = false,
  });

  @override
  Widget build(BuildContext context) {
    return Card(
      margin: const EdgeInsets.symmetric(vertical: 8),
      color: Colors.grey[850],
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      child: ListTile(
        leading: CircleAvatar(
          backgroundColor: Colors.tealAccent,
          child: Text(company[0]),
        ),
        title: Text(title, style: const TextStyle(fontWeight: FontWeight.bold)),
        subtitle: Text("$company • $location"),
        trailing: Icon(
          isSaved ? Icons.star : Icons.star_border,
          color: Colors.tealAccent,
        ),
        onTap: () {
          // expand / collapse details within the card
	  showModalBottomSheet(
	    context: context,
	    backgroundColor: Colors.grey[900],
	    shape: const RoundedRectangleBorder(
	      borderRadius: BorderRadius.vertical(top: Radius.circular(16)),
	    ),
	    builder: (context) {
	      return Padding(
	        padding: const EdgeInsets.all(16.0),
	        child: Column(
	          mainAxisSize: MainAxisSize.min,
	          crossAxisAlignment: CrossAxisAlignment.start,
	          children: [
	            Text(title, style: const TextStyle(fontSize: 24, fontWeight: FontWeight.bold)),
	            const SizedBox(height: 8),
	            Text("$company • $location", style: const TextStyle(fontSize: 16, color: Colors.white70)),
	            const SizedBox(height: 16),
	            const Text(
	              "Job Description",
	              style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
	            ),
	            const SizedBox(height: 8),
	            const Text(
	              "This is a placeholder for the job description. Here you would find details about the job responsibilities, requirements, and other relevant information.",
	              style: TextStyle(fontSize: 16),
	            ),
	            const SizedBox(height: 16),
	            ElevatedButton.icon(
	              onPressed: () {
	                Navigator.pop(context);
	              },
	              icon: Icon(isSaved ? Icons.star : Icons.star_border),
	              label: Text(isSaved ? "Unsave Job" : "Save Job"),
	              style: ElevatedButton.styleFrom(
	                backgroundColor: Colors.tealAccent,
	                foregroundColor: Colors.black87,
	                minimumSize: const Size.fromHeight(50),
	              ),
	            ),
	          ],
	        ),
	      );
	    },

	  );
        },
      ),
    );
  }
}

