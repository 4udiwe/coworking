import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:coworking_app/core/l10n/localization_helper.dart';

import '../../../../core/models/coworking.dart';
import '../../bloc/coworking_bloc.dart';
import '../../bloc/coworking_event.dart';
import '../screens/coworking_details_screen.dart';

class CoworkingCard extends StatelessWidget {
  final Coworking coworking;

  const CoworkingCard({super.key, required this.coworking});

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: () {
        if (!coworking.isActive) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              duration: const Duration(seconds: 1),
              content: Text(context.l10n.coworkingUnavailable),
            ),
          );
          return;
        }
        final bloc = context.read<CoworkingBloc>();
        bloc.add(SelectCoworking(coworking.id));

        Navigator.push(
          context,
          MaterialPageRoute(
            builder: (_) => BlocProvider.value(
              value: bloc,
              child: CoworkingDetailsScreen(coworkingId: coworking.id),
            ),
          ),
        );
      },
      child: Card(
        shadowColor: Colors.blue,
        surfaceTintColor: coworking.isActive ? Colors.white : Colors.black,
        elevation: 12,
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            /// 🔷 Image (заглушка)
            Container(
              height: 250,
              decoration: BoxDecoration(
                color: Colors.blue.withOpacity(0.4),
                borderRadius: const BorderRadius.vertical(
                  top: Radius.circular(12),
                ),
              ),
              child: Stack(
                children: [
                  const Center(child: Icon(Icons.business, size: 100)),
                  Center(
                    child: !coworking.isActive
                        ? Icon(
                            Icons.block,
                            size: 200,
                            color: Colors.grey.withAlpha(150),
                          )
                        : null,
                  ),
                ],
              ),
            ),

            /// 🔷 Info
            Padding(
              padding: const EdgeInsets.all(12),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    coworking.name,
                    style: const TextStyle(
                      fontWeight: FontWeight.bold,
                      fontSize: 16,
                    ),
                  ),
                  const SizedBox(height: 4),
                  Text(
                    coworking.address,
                    style: TextStyle(color: Colors.grey[600], fontSize: 12),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}
