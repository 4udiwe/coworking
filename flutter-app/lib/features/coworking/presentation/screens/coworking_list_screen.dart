import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:coworking_app/core/l10n/localization_helper.dart';
import '../../../../core/di/service_locator.dart';
import '../../../../core/models/coworking.dart';
import '../../../../core/utils/bloc_load_state.dart';
import '../../bloc/coworking_bloc.dart';
import '../../bloc/coworking_event.dart';
import '../../bloc/coworking_state.dart';
import '../widgets/coworking_card.dart';

class CoworkingPage extends StatelessWidget {
  const CoworkingPage({super.key});

  @override
  Widget build(BuildContext context) {
    return BlocProvider(
      create: (_) => sl<CoworkingBloc>()..add(FetchCoworkings()),
      child: const CoworkingListScreen(),
    );
  }
}

class CoworkingListScreen extends StatefulWidget {
  const CoworkingListScreen({super.key});

  @override
  State<CoworkingListScreen> createState() => _CoworkingListScreenState();
}

class _CoworkingListScreenState extends State<CoworkingListScreen> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(context.l10n.coworkingsScreenTitle),
        automaticallyImplyLeading: false,
      ),
      body: BlocBuilder<CoworkingBloc, CoworkingState>(
        builder: (context, state) {
          final coworkingsState = state.coworkings;

          switch (coworkingsState.status) {
            case LoadStatus.loading:
              return const Center(child: CircularProgressIndicator());
            case LoadStatus.error:
              return Center(
                child: Text(
                  coworkingsState.error ?? "Ошибка загрузки коворкингов",
                ),
              );
            case LoadStatus.success:
              final coworkings = coworkingsState.data ?? [];
              if (coworkings.isEmpty) {
                return const Center(child: Text("Нет коворкингов"));
              }

              return LayoutBuilder(
                builder: (context, constraints) {
                  final isWide = constraints.maxWidth > 800;

                  if (isWide) {
                    return _Grid(coworkings);
                  } else {
                    return _List(coworkings);
                  }
                },
              );
            default:
              return const SizedBox();
          }
        },
      ),
    );
  }
}

class _List extends StatelessWidget {
  final List<Coworking> coworkings;
  const _List(this.coworkings);

  @override
  Widget build(BuildContext context) {
    return ListView.builder(
      itemCount: coworkings.length,
      itemBuilder: (context, index) {
        final coworking = coworkings[index];
        return CoworkingCard(coworking: coworking);
      },
    );
  }
}

class _Grid extends StatelessWidget {
  final List<Coworking> coworkings;
  const _Grid(this.coworkings);

  @override
  Widget build(BuildContext context) {
    return GridView.builder(
      padding: const EdgeInsets.all(16),
      gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
        crossAxisCount: 2,
        crossAxisSpacing: 16,
        mainAxisSpacing: 16,
        childAspectRatio: 2,
      ),
      itemCount: coworkings.length,
      itemBuilder: (context, index) {
        final coworking = coworkings[index];
        return CoworkingCard(coworking: coworking);
      },
    );
  }
}
